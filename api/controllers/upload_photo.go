package controllers

import (
	"backend/api/httputil"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"firebase.google.com/go/auth"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// UploadPhoto maneja la solicitud para cargar una foto de perfil de usuario.
//
// @Summary Cargar foto de perfil
// @Description Cargar una foto de perfil para el usuario actual.
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Token de autorización JWT"
// @Param file formData file true "Archivo de imagen"
// @Success 200 {object} httputil.StandardResponse "Foto de perfil cargada correctamente"
// @Failure 400 {object} httputil.ErrorResponse "No se ha enviado ningún archivo"
// @Failure 401 {object} httputil.ErrorResponse "No autorizado"
// @Failure 500 {object} httputil.ErrorResponse "Error al cargar la foto de perfil"
// @Router /upload-photo [post]
func UploadPhoto(c *gin.Context, firestoreClient *firestore.Client, storageClient *storage.Client, authClient *auth.Client) {
	// Verificar si se ha enviado un archivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "No se ha enviado ningún archivo"})
		return
	}

	// Abrir el archivo temporalmente
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al abrir el archivo"})
		return
	}
	defer src.Close()

	// Verificar el tipo MIME del archivo
	mime, err := mimetype.DetectReader(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al detectar el tipo MIME del archivo"})
		return
	}

	if !mime.Is("image/jpeg") && !mime.Is("image/png") && !mime.Is("image/gif") && !mime.Is("image/webp") {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "El archivo debe ser una imagen JPEG, PNG, GIF o WebP"})
		return
	}

	// Volver al inicio del archivo para la siguiente lectura
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al volver al inicio del archivo"})
		return
	}

	// Redimensionar la imagen
	img, err := imaging.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al decodificar la imagen"})
		return
	}

	resizedImg := imaging.Resize(img, 200, 200, imaging.Lanczos)

	// Crear un archivo temporal para la imagen redimensionada
	tempFilePath := filepath.Join("./uploads", uuid.New().String()+filepath.Ext(file.Filename))
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al crear el archivo temporal"})
		return
	}

	// Asegurarse de cerrar el archivo al final de la función, antes del retorno
	defer func() {
		tempFile.Close() // Cerrar el archivo temporal

		// Eliminar el archivo temporal después de cargarlo en la nube y cerrar el archivo
		if err := os.Remove(tempFilePath); err != nil {
			log.Printf("Advertencia: no se pudo eliminar el archivo temporal %s: %v", tempFilePath, err)
		}
	}()

	// Guardar la imagen redimensionada en el archivo temporal
	err = imaging.Encode(tempFile, resizedImg, imaging.JPEG)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al guardar la imagen redimensionada"})
		return
	}

	// Obtener el token de usuario del contexto
	token, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "No autorizado"})
		return
	}

	// Extraer UID del token
	uid := token.(*auth.Token).UID

	// Eliminar todas las fotos anteriores del usuario en Cloud Storage
	if err := deleteAllFromCloudStorage(uid, storageClient); err != nil {
		log.Printf("Advertencia: no se pudieron eliminar las fotos anteriores de Cloud Storage: %v", err)
	}

	// Subir el archivo al servicio de almacenamiento en la nube (ejemplo: Google Cloud Storage)
	url, err := uploadToCloudStorage(uid, filepath.Base(tempFilePath), tempFilePath, storageClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar la foto de perfil"})
		return
	}

	// Actualizar el perfil del usuario en Firestore con la URL de la foto de perfil
	_, err = firestoreClient.Collection("users").Doc(uid).Set(context.Background(), map[string]interface{}{
		"photoURL": url,
	}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar la foto de perfil"})
		return
	}
	// Actualizar el perfil del usuario en Firebase Authentication con la URL de la foto de perfil
	_, err = authClient.UpdateUser(context.Background(), uid, (&auth.UserToUpdate{}).PhotoURL(url))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al actualizar la foto de perfil"})
		return
	} else {
		log.Printf("Foto de perfil actualizada para el usuario %s", uid)
	}

	// Crear un token renovado con la URL de la foto de perfil
	// para que el cliente pueda actualizar la foto de perfil en la caché
	// sin necesidad de volver a iniciar sesión
	userClaims, error := GetUserInfo(authClient, uid)
	if error != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al obtener información del usuario"})
		return
	}

	c.JSON(http.StatusOK, httputil.StandardResponse{
		Message: "Foto de perfil cargada correctamente",
		Data:    userClaims,
	})
}

// deleteAllFromCloudStorage elimina todos los archivos de la carpeta del usuario en el almacenamiento en la nube.
func deleteAllFromCloudStorage(uid string, client *storage.Client) error {
	// Configurar contexto y cliente para Google Cloud Storage
	ctx := context.Background()

	// Obtener el bucket del entorno
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return fmt.Errorf("GCS_BUCKET_NAME no está configurado en las variables de entorno")
	}

	// Construir el prefijo del objeto en el bucket de almacenamiento
	prefix := fmt.Sprintf("profile_photos/%s/", uid)

	// Obtener todos los objetos en el bucket con el prefijo dado
	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{
		Prefix: prefix,
	})

	// Iterar sobre los objetos y eliminar cada uno
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error al obtener objetos en Cloud Storage: %v", err)
		}

		// Eliminar el objeto del bucket
		if err := client.Bucket(bucketName).Object(objAttrs.Name).Delete(ctx); err != nil {
			return fmt.Errorf("error al eliminar objeto en Cloud Storage: %v", err)
		}
	}

	return nil
}

// uploadToCloudStorage sube un archivo al almacenamiento en la nube (ejemplo: Google Cloud Storage)
func uploadToCloudStorage(uid, fileName, filePath string, client *storage.Client) (string, error) {
	// Configurar contexto y cliente para Google Cloud Storage
	ctx := context.Background()

	// Abrir el archivo local para lectura
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo local %s: %v", filePath, err)
	}
	defer file.Close()

	// Obtener el bucket del entorno
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("GCS_BUCKET_NAME no está configurado en las variables de entorno")
	}

	// Construir el nombre del objeto en el bucket de almacenamiento
	objName := fmt.Sprintf("profile_photos/%s/%s", uid, fileName)

	// Obtener el manejador del objeto en el bucket
	wc := client.Bucket(bucketName).Object(objName).NewWriter(ctx)

	// Copiar el contenido del archivo al objeto en el bucket
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("error al copiar archivo a Cloud Storage: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error al cerrar escritor de Cloud Storage: %v", err)
	}

	// Configurar el objeto para ser público
	acl := client.Bucket(bucketName).Object(objName).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("error al hacer público el objeto: %v", err)
	}

	// Obtener la URL pública del objeto cargado
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objName)

	return url, nil
}

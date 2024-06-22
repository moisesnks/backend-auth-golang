package controllers

import (
	emailPkg "backend/api/email"
	"backend/api/httputil"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	// Cargar el template de correo de verificación
	verificationT = template.Must(template.ParseFiles("html/verification_email.html"))
}

// RegisterUser maneja el registro de usuario utilizando la API de Firebase Authentication.
//
// Este endpoint registra un nuevo usuario utilizando Firebase Authentication.
// Se espera recibir datos JSON que contengan el correo electrónico y la contraseña del usuario.
// Después de registrar al usuario exitosamente en Firebase Auth, se guarda la información del usuario
// en Firestore, se genera un código de verificación y se envía un correo de verificación al usuario.
//
// @Summary Registrar usuario
// @Description Registra un nuevo usuario utilizando Firebase Authentication y envía un correo de verificación.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Datos de registro del usuario"
// @Success 200 {object} httputil.StandardResponse "Respuesta exitosa al registrar usuario"
// @Failure 400 {object} httputil.ErrorResponse "Datos de solicitud inválidos"
// @Failure 401 {object} httputil.ErrorResponse "El correo electrónico ya está en uso"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /register [post]
func RegisterUser(c *gin.Context, firestoreClient *firestore.Client) {
	// validar firestoreClient
	if firestoreClient == nil {
		log.Println("Firestore client no inicializado")
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
		return
	}

	var registerData RegisterRequest
	if err := c.BindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Datos de solicitud inválidos"})
		return
	}

	requestData := map[string]interface{}{
		"email":             registerData.Email,
		"password":          registerData.Password,
		"returnSecureToken": true,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
		return
	}

	firebaseAPIURL := "https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=" + firebaseAPIKey

	resp, err := http.Post(firebaseAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error al realizar la solicitud a Firebase: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al comunicarse con Firebase"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error al decodificar la respuesta de Firebase: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al procesar la respuesta de Firebase"})
		return
	}

	if errMsg, ok := result["error"].(map[string]interface{}); ok {
		errorMessage := errMsg["message"].(string)
		switch errorMessage {
		case "EMAIL_EXISTS":
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "El correo electrónico ya está en uso"})
		default:
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: errorMessage})
		}
		return
	}

	// Usuario creado exitosamente
	uid := result["localId"].(string)
	email := result["email"].(string)
	verificationCode := generateVerificationCode() // Envío de correo de verificación

	// Guardar datos en Firestore
	docRef := firestoreClient.Collection("users").Doc(uid)
	userData := map[string]interface{}{
		"email":            email,
		"createdAt":        time.Now(),
		"verified":         false,
		"verificationCode": verificationCode,
		"codeValidUntil":   time.Now().Add(30 * time.Minute),
	}
	_, err = docRef.Set(context.Background(), userData)
	if err != nil {
		log.Printf("Error al guardar datos en Firestore: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al guardar datos en Firestore"})
		return
	}

	// Configurar y enviar el correo de verificación
	mailSubject := "Código de verificación"

	// Datos para el template HTML del correo
	data := struct {
		VerificationCode string
		Email            string
		UID              string
		Year             int
	}{
		VerificationCode: verificationCode,
		Email:            email,
		UID:              uid,
		Year:             time.Now().Year(),
	}

	var mailBody bytes.Buffer
	if err := verificationT.Execute(&mailBody, data); err != nil {
		log.Printf("Error al ejecutar el template HTML: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al procesar el template HTML"})
		return
	}

	// Utiliza la función SendMail del paquete email
	smtpConfig := emailPkg.SmtpConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	err = emailPkg.SendMail(smtpConfig, email, mailSubject, mailBody.String())
	if err != nil {
		log.Printf("Error al enviar el correo de verificación: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al enviar el correo de verificación"})
		return
	}

	response := httputil.StandardResponse{
		Message: "Usuario registrado exitosamente. Se ha enviado un correo de verificación.",
		Data:    map[string]string{"uid": uid, "email": email},
	}
	c.JSON(http.StatusOK, response)
}

// Función para generar un código de verificación (simulado)
func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

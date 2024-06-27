// updateProfile.go

package controllers

import (
	"context"
	"log"
	"net/http"

	"backend/api/httputil"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// Estructura para los datos de actualización de perfil
type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName"`
	RUT         *string `json:"rut"`
	Birthdate   *string `json:"birthdate"`
}

// Función para actualizar el perfil del usuario
// UpdateProfile actualiza el perfil de un usuario en Firebase Authentication y Firestore.
//
// Esta función permite actualizar el displayName en Firebase Auth y los campos rut y birthdate en Firestore
// para el usuario autenticado.
//
// @Summary Actualiza el perfil de usuario
// @Description Actualiza el displayName en Firebase Auth y los campos rut y birthdate en Firestore
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización JWT"
// @Param body body UpdateProfileRequest true "Datos de actualización del perfil"
// @Success 200 {object} httputil.StandardResponse "Perfil actualizado exitosamente"
// @Failure 400 {object} httputil.ErrorResponse "Datos de solicitud inválidos"
// @Failure 401 {object} httputil.ErrorResponse "Usuario no autorizado"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /update [patch]
// Función para actualizar el perfil del usuario
func UpdateProfile(c *gin.Context, firestoreClient *firestore.Client, authClient *auth.Client) {
	var req UpdateProfileRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Datos de solicitud inválidos"})
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

	// Construir parámetros para actualizar en Firebase Auth
	params := (&auth.UserToUpdate{})

	// Actualizar el displayName si se proporciona
	if req.DisplayName != nil {
		params = params.DisplayName(*req.DisplayName)
	}

	// Actualizar en Firebase Auth si se proporciona un displayName
	var u *auth.UserRecord
	var err error
	if req.DisplayName != nil {
		u, err = authClient.UpdateUser(context.Background(), uid, params)
		if err != nil {
			log.Printf("Error actualizando usuario en Firebase Auth: %v", err)
			c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
			return
		}
		// Actualizar el displayName en Firestore si se proporciona
		userRef := firestoreClient.Collection("users").Doc(uid)
		_, err := userRef.Set(context.Background(), map[string]interface{}{"displayName": *req.DisplayName}, firestore.MergeAll)
		if err != nil {
			log.Printf("Error actualizando usuario en Firestore: %v", err)
			c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
			return
		}
	}

	// Actualizar rut y/o birthdate en Firestore si se proporciona alguno
	if req.RUT != nil || req.Birthdate != nil {
		updateData := map[string]interface{}{}

		if req.RUT != nil {
			updateData["rut"] = *req.RUT
		}
		if req.Birthdate != nil {
			updateData["birthdate"] = *req.Birthdate
		}

		userRef := firestoreClient.Collection("users").Doc(uid)
		_, err := userRef.Set(context.Background(), updateData, firestore.MergeAll)
		if err != nil {
			log.Printf("Error actualizando usuario en Firestore: %v", err)
			c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
			return
		}
	}

	// Preparar la respuesta JSON
	var response httputil.StandardResponse
	if u != nil {
		data := map[string]interface{}{
			"displayName": u.DisplayName,
			"email":       u.Email,
			"uid":         uid,
		}

		// Agregar rut si se proporcionó
		if req.RUT != nil {
			data["rut"] = *req.RUT
		}

		// Agregar birthdate si se proporcionó
		if req.Birthdate != nil {
			data["birthdate"] = *req.Birthdate
		}

		response = httputil.StandardResponse{
			Message: "Perfil actualizado exitosamente",
			Data:    data,
		}
	} else {
		response = httputil.StandardResponse{
			Message: "Perfil actualizado exitosamente",
		}
	}

	// Retornar la respuesta al cliente
	c.JSON(http.StatusOK, response)
}

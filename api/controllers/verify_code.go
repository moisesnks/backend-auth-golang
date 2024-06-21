package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type VerifyCodeRequest struct {
	UID              string `json:"uid"`
	VerificationCode string `json:"verificationCode"`
}

type VerifyCodeResponse struct {
	Message string `json:"message"`
}

// VerifyCode verifica el código de verificación de un usuario.
//
// Este endpoint verifica si el código proporcionado por el usuario coincide con el
// código almacenado en Firestore para el usuario identificado por UID. Si el código
// es correcto y aún no se ha verificado, marca al usuario como verificado y asigna
// el rol de "miembro" como un custom claim en Firebase Auth.
//
// @Summary Verifica el código de verificación de un usuario.
// @Description Verifica el código de verificación de un usuario en Firestore y Firebase Auth.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body VerifyCodeRequest true "Datos de la solicitud"
// @Success 200 {object} VerifyCodeResponse "Respuesta exitosa al verificar el código"
// @Failure 400 {object} ErrorResponse "Datos de solicitud inválidos"
// @Failure 401 {object} ErrorResponse "Código de verificación incorrecto"
// @Failure 404 {object} ErrorResponse "Usuario no encontrado"
// @Failure 500 {object} ErrorResponse "Error interno del servidor"
// @Router /verify [post]
func VerifyCode(c *gin.Context, firestoreClient *firestore.Client, authClient *auth.Client) {
	var req VerifyCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Datos de solicitud inválidos"})
		return
	}

	// acceder a la colección de usuarios
	users := firestoreClient.Collection("users")

	// verificar si el usuario existe
	doc, err := users.Doc(req.UID).Get(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error interno del servidor"})
		return
	}
	if !doc.Exists() {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Usuario no encontrado"})
		return
	}

	now := time.Now()

	// verificar que el codeValidUntil sea mayor a la fecha actual
	userData := doc.Data()
	codeValidUntil, ok := userData["codeValidUntil"].(time.Time)
	if !ok || codeValidUntil.Before(now) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Código de verificación expirado"})
		return
	}

	// verificar que el código de verificación sea correcto
	verificationCode, ok := userData["verificationCode"].(string)
	if !ok || verificationCode != req.VerificationCode {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Código de verificación incorrecto"})
		return
	}

	// verificar que el usuario no haya sido verificado previamente
	verified, ok := userData["verified"].(bool)
	if !ok || verified {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Usuario ya verificado"})
		return
	}

	// marcar el usuario como verificado
	_, err = users.Doc(req.UID).Set(c, map[string]interface{}{
		"verified": true,
	}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error interno del servidor"})
		return
	}

	// Asignar el rol de miembro como custom claim
	err = authClient.SetCustomUserClaims(context.Background(), req.UID, map[string]interface{}{"role": "member"})
	if err != nil {
		log.Printf("Error al asignar rol de miembro: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario verificado correctamente"})
}

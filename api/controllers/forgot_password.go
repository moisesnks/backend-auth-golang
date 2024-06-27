package controllers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/api/email"
	"backend/api/httputil"
	"backend/api/utils"
)

var resetPasswordT *template.Template

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}
	resetPasswordT = template.Must(template.ParseFiles("html/reset_password_email.html"))
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// ForgotPassword maneja las solicitudes de restablecimiento de contraseña.
//
// Este endpoint recibe un correo electrónico, verifica si el usuario existe en Firebase Authentication
// y envía un correo electrónico con un enlace para restablecer la contraseña.
//
// @Summary Restablecer contraseña
// @Description Envía un correo electrónico para restablecer la contraseña si el usuario existe en Firebase Authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body ForgotPasswordRequest true "Correo electrónico del usuario"
// @Success 200 {object} httputil.StandardResponse "Correo de restablecimiento de contraseña enviado"
// @Failure 400 {object} httputil.ErrorResponse "Datos de solicitud inválidos"
// @Failure 404 {object} httputil.ErrorResponse "Usuario no encontrado"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /forgot-password [post]
func ForgotPassword(c *gin.Context, authClient *auth.Client, firestoreClient *firestore.Client) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Datos de solicitud inválidos"})
		return
	}

	user, err := authClient.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse{Message: "Usuario no encontrado"})
		return
	}

	// Generar token JWT
	token, err := utils.GenerarToken(req.Email, user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al generar token de restablecimiento"})
		return
	}

	// Guardar token en Firestore con una expiración
	resetData := map[string]interface{}{
		"email":      req.Email,
		"resetToken": token,
		"expiresAt":  time.Now().Add(30 * time.Minute), // Establecer expiración del token si es necesario
	}
	_, err = firestoreClient.Collection("password_resets").Doc(req.Email).Set(context.Background(), resetData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al guardar token de restablecimiento"})
		return
	}

	// Construir el enlace de restablecimiento de contraseña
	urlFrontend := os.Getenv("URL_FRONTEND")
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", urlFrontend, token)
	subject := "Restablecimiento de contraseña"

	// Datos para el template HTML del correo
	data := struct {
		ResetLink string
		Email     string
		Year      int
	}{
		ResetLink: resetLink,
		Email:     req.Email,
		Year:      time.Now().Year(),
	}

	var mailBody bytes.Buffer
	if err := resetPasswordT.Execute(&mailBody, data); err != nil {
		log.Printf("Error al ejecutar el template HTML: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al enviar el correo de restablecimiento"})
		return
	}

	// Configurar y enviar el correo de restablecimiento
	smtpConfig := email.SmtpConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	err = email.SendMail(smtpConfig, req.Email, subject, mailBody.String())
	if err != nil {
		log.Printf("Error al enviar el correo de restablecimiento: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al enviar el correo de restablecimiento"})
		return
	}

	// Enviar respuesta de éxito
	response := httputil.StandardResponse{
		Message: "Se ha enviado un correo electrónico con instrucciones para restablecer la contraseña.",
	}
	c.JSON(http.StatusOK, response)
}

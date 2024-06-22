// backend/api/controllers/resend_code.go

package controllers

import (
	emailPkg "backend/api/email"
	"backend/api/httputil"
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func init() {
	// Cargar el template de correo de verificación
	verificationT = template.Must(template.ParseFiles("html/verification_email.html"))
}

// ResendCode reenvía el código de verificación si es necesario.
//
// Este endpoint permite reenviar el código de verificación a un usuario no verificado cuyo
// código anterior haya expirado. Se espera que el usuario esté autenticado mediante un token
// de sesión válido en el encabezado de la solicitud.
//
// @Summary Reenviar código de verificación
// @Description Reenvía el código de verificación al usuario si el código anterior ha expirado.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httputil.StandardResponse "Respuesta exitosa al reenviar código de verificación"
// @Failure 400 {object} httputil.ErrorResponse "El usuario ya está verificado o el código de verificación aún es válido"
// @Failure 401 {object} httputil.ErrorResponse "Usuario no autorizado"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /resend-code [post]
func ResendCode(c *gin.Context, firestoreClient *firestore.Client) {
	// Obtener el token de usuario del contexto
	token, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "No autorizado"})
		return
	}

	// Extraer UID del token
	uid := token.(*auth.Token).UID

	// Obtener datos del usuario desde Firestore
	userData, err := GetUserData(uid, firestoreClient)
	if err != nil {
		log.Printf("Error al obtener datos del usuario desde Firestore: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al obtener datos del usuario"})
		return
	}

	// Verificar si el usuario ya está verificado
	if userData.Verified {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "El usuario ya está verificado"})
		return
	}

	// Verificar si el código de verificación aún es válido
	if IsVerificationCodeValid(userData.CodeValidUntil) {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "El código de verificación aún es válido"})
		return
	}

	// Generar un nuevo código de verificación y actualizar los datos del usuario en Firestore
	newVerificationCode := generateVerificationCode()
	userData.VerificationCode = newVerificationCode
	userData.CodeValidUntil = time.Now().Add(30 * time.Minute)

	_, err = firestoreClient.Collection("users").Doc(uid).Set(context.Background(), userData)
	if err != nil {
		log.Printf("Error al actualizar datos del usuario en Firestore: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al actualizar datos del usuario"})
		return
	}

	// Configurar y enviar el correo de verificación
	mailSubject := "Código de verificación"

	data := struct {
		VerificationCode string
		Email            string
		UID              string
		Year             int
	}{
		VerificationCode: newVerificationCode,
		Email:            userData.Email,
		UID:              uid,
		Year:             time.Now().Year(),
	}

	var mailBody bytes.Buffer
	if err := verificationT.Execute(&mailBody, data); err != nil {
		log.Printf("Error al ejecutar el template HTML: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al procesar el template HTML"})
		return
	}

	smtpConfig := emailPkg.SmtpConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	err = emailPkg.SendMail(smtpConfig, userData.Email, mailSubject, mailBody.String())
	if err != nil {
		log.Printf("Error al enviar el correo de verificación: %v", err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al enviar el correo de verificación"})
		return
	}

	// Respuesta exitosa
	response := httputil.StandardResponse{
		Message: "Se ha enviado un nuevo correo de verificación.",
	}
	c.JSON(http.StatusOK, response)
}

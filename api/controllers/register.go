package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// RegisterResponse representa la estructura de la respuesta exitosa
type RegisterResponse struct {
	Message string `json:"message"`
	UID     string `json:"uid"`
	Email   string `json:"email"`
}

var (
	smtpConfig    SmtpConfig // Configuración SMTP para el envío de correos
	verificationT *template.Template
)

type SmtpConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func init() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}
	// Configuración SMTP para el envío de correos
	smtpConfig = SmtpConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	// Cargar el template de correo de verificación
	verificationT = template.Must(template.ParseFiles("html/verification_email.html"))

}

// RegisterData representa los datos para el registro de usuario
type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
// @Param body body RegisterData true "Datos de registro del usuario"
// @Success 200 {object} RegisterResponse "Respuesta exitosa al registrar usuario"
// @Failure 400 {object} ErrorResponse "Datos de solicitud inválidos"
// @Failure 401 {object} ErrorResponse "El correo electrónico ya está en uso"
// @Failure 500 {object} ErrorResponse "Error interno del servidor"
// @Router /register [post]
func RegisterUser(c *gin.Context, firestoreClient *firestore.Client) {
	// validar firestoreClient
	if firestoreClient == nil {
		log.Println("Firestore client no inicializado")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	var registerData RegisterData
	if err := c.BindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de solicitud inválidos"})
		return
	}

	requestData := map[string]interface{}{
		"email":             registerData.Email,
		"password":          registerData.Password,
		"returnSecureToken": true,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	firebaseAPIURL := "https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=" + firebaseAPIKey

	resp, err := http.Post(firebaseAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error al realizar la solicitud a Firebase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al comunicarse con Firebase"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error al decodificar la respuesta de Firebase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta de Firebase"})
		return
	}

	if errMsg, ok := result["error"].(map[string]interface{}); ok {
		errorMessage := errMsg["message"].(string)
		switch errorMessage {
		case "EMAIL_EXISTS":
			c.JSON(http.StatusBadRequest, gin.H{"error": "El correo electrónico ya está en uso"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar datos en Firestore"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el template HTML"})
		return
	}

	err = sendMail(email, mailSubject, mailBody.String())
	if err != nil {
		log.Printf("Error al enviar el correo de verificación: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo de verificación"})
		return
	}

	response := RegisterResponse{
		Message: "Usuario creado exitosamente y correo de verificación enviado",
		UID:     uid,
		Email:   email,
	}
	c.JSON(http.StatusOK, response)
}

// Función para generar un código de verificación (simulado)
func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Función para enviar correo electrónico usando SMTP
func sendMail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, smtpConfig.Username, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/api/httputil"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Variables globales
var (
	firebaseAPIKey string
)

// Función de inicialización para cargar la API_KEY desde el archivo .env
func init() {
	// Load Firebase API key from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	firebaseAPIKey = os.Getenv("FIREBASE_API_KEY")
	if firebaseAPIKey == "" {
		log.Fatal("FIREBASE_API_KEY is not set in .env file")
	}

	fmt.Println("Firebase API Key: ", firebaseAPIKey)

}

// LoginRequest representa los datos de inicio de sesión del usuario
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser maneja el inicio de sesión de usuario utilizando la API de Firebase Authentication
// @Summary Inicia sesión de usuario
// @Description Inicia sesión utilizando Firebase Authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param email body LoginRequest true "Datos de inicio de sesión"
// @Success 200 {object} httputil.StandardResponse "Inicio de sesión exitoso"
// @Failure 400 {object} httputil.ErrorResponse "Datos de solicitud inválidos o errores en la solicitud"
// @Failure 401 {object} httputil.ErrorResponse "Credenciales incorrectas"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var loginData LoginRequest
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Datos de solicitud inválidos"})
		return
	}

	requestData := map[string]string{
		"email":             loginData.Email,
		"password":          loginData.Password,
		"returnSecureToken": "true",
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error interno del servidor"})
		return
	}

	firebaseAPIURL := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + firebaseAPIKey

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
		case "EMAIL_NOT_FOUND":
			c.JSON(http.StatusNotFound, httputil.ErrorResponse{Message: "Usuario no encontrado"})
		case "INVALID_PASSWORD":
			c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Contraseña incorrecta"})
		default:
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: errorMessage})
		}
		return
	}

	response := httputil.StandardResponse{
		Message: "Inicio de sesión exitoso",
		Data:    map[string]string{"token": result["idToken"].(string)},
	}
	c.JSON(http.StatusOK, response)
}

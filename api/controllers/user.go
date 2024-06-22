// backend/api/controllers/user.go

package controllers

import (
	"context"
	"log"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	verificationT *template.Template
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

// RegisterRequest representa los datos para el registro de usuario
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserData representa los datos del usuario en Firestore
type UserData struct {
	Email            string    `firestore:"email"`
	Verified         bool      `firestore:"verified"`
	VerificationCode string    `firestore:"verificationCode"`
	CodeValidUntil   time.Time `firestore:"codeValidUntil"`
}

// Función auxiliar para obtener los datos del usuario desde Firestore
func GetUserData(uid string, firestoreClient *firestore.Client) (UserData, error) {
	var userData UserData

	docRef := firestoreClient.Collection("users").Doc(uid)
	snapshot, err := docRef.Get(context.Background())
	if err != nil {
		return userData, err
	}

	if err := snapshot.DataTo(&userData); err != nil {
		return userData, err
	}

	return userData, nil
}

// Función auxiliar para verificar si el código de verificación sigue válido
func IsVerificationCodeValid(codeValidUntil time.Time) bool {
	return time.Now().Before(codeValidUntil)
}

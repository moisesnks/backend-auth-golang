package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte(os.Getenv("SECRET_KEY"))

// Claims estructura para almacenar los claims del token JWT
type Claims struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
	jwt.StandardClaims
}

// GenerarToken genera un JWT para restablecimiento de contraseña
func GenerarToken(email, uid string) (string, error) {
	// Crear token con un payload (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		email,
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(), // Token expira en 30 minutos
		},
	})

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %v", err)
	}

	return tokenString, nil
}

// VerificarToken verifica un token JWT y devuelve los claims si es válido
func VerificarToken(tokenString string) (*Claims, error) {
	// Parsear el token con la clave secreta y extraer los claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error al parsear el token: %v", err)
	}

	// Verificar si el token es válido
	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Obtener los claims del token
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("error al obtener los claims del token")
	}

	return claims, nil
}

// IsTokenExpired verifica si el token JWT ha expirado
func (c *Claims) IsTokenExpired() bool {
	return c.ExpiresAt < time.Now().Unix()
}

// middleware/auth.go

package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el token de autorización JWT y coloca el usuario en el contexto de Gin si el token es válido.
func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización no proporcionado"})
			c.Abort()
			return
		}

		idToken := extractTokenFromHeader(authHeader)
		if idToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token de autorización inválido"})
			c.Abort()
			return
		}

		// Verificar el token JWT
		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización inválido"})
			c.Abort()
			return
		}

		// Colocar el usuario en el contexto
		c.Set("user", token)
		c.Next()
	}
}

// Función para extraer el token de autorización del encabezado
func extractTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

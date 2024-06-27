package middleware

import (
	"backend/api/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware es un middleware para verificar y decodificar el token JWT del header Authorization
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización no proporcionado"})
			c.Abort()
			return
		}

		// Verificar el esquema del token (debe ser Bearer)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Obtener el token JWT
		tokenString := parts[1]

		// Verificar y decodificar el token
		claims, err := utils.VerificarToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Token inválido: %v", err)})
			c.Abort()
			return
		}

		// Pasar los claims al contexto para que el siguiente controlador pueda usarlos
		c.Set("claims", claims)

		// Llamar al siguiente controlador
		c.Next()
	}
}

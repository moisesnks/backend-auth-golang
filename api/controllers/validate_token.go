package controllers

import (
	"backend/api/httputil"
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// ValidateToken devuelve la información que hay en el token que ya pasó el middleware de autenticación
//
// Este endpoint permite validar un token de usuario y obtener la información que contiene.
// Se espera que el usuario esté autenticado mediante un token de sesión válido en el encabezado de la solicitud.
//
// @Summary Validar token de usuario
// @Description Valida un token de usuario y devuelve la información que contiene
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httputil.StandardResponse "Respuesta exitosa al validar token de usuario"
// @Failure 401 {object} httputil.ErrorResponse "Usuario no autorizado"
// @Router /validate-token [get]
// ValidateToken verifica y devuelve la información del token validado.
func ValidateToken(c *gin.Context, authClient *auth.Client) {
	// Obtener el token de usuario del contexto
	token, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "No autorizado"})
		return
	}

	// Convertir token a *auth.Token
	authTok, ok := token.(*auth.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al obtener token"})
		return
	}

	// Construir el objeto de usuario desde el token decodificado
	user := gin.H{
		"uid":           authTok.UID,
		"email":         authTok.Claims["email"],
		"role":          authTok.Claims["role"],
		"emailVerified": authTok.Claims["email_verified"],
		"displayName":   authTok.Claims["name"],
		"photoURL":      authTok.Claims["picture"],
	}

	// Enviar la respuesta con el usuario válido
	c.JSON(http.StatusOK, httputil.StandardResponse{
		Message: "Token válido",
		Data:    user,
	})
}

// la misma función pero útil para que otras funciones la llamen
func GetUserInfo(authClient *auth.Client, uid string) (map[string]interface{}, error) {
	// Obtener el usuario desde Firebase Authentication
	user, err := authClient.GetUser(context.Background(), uid)
	if err != nil {
		return nil, err
	}

	// Construir el objeto de usuario
	userInfo := map[string]interface{}{
		"uid":           user.UID,
		"email":         user.Email,
		"role":          user.CustomClaims["role"],
		"emailVerified": user.EmailVerified,
		"displayName":   user.DisplayName,
		"photoURL":      user.PhotoURL,
	}

	return userInfo, nil
}

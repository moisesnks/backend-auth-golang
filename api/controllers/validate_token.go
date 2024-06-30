package controllers

import (
	"backend/api/httputil"
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

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

	// Obtener la información del usuario desde Firebase Authentication usando el uid
	userInfo, err := GetUserInfo(authClient, authTok.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al obtener información del usuario"})
		return
	}

	// Actualizar los claims del token si es necesario
	if needToUpdateClaims(authTok.Claims, userInfo) {
		err := updateClaims(authClient, authTok.UID, userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al actualizar los claims"})
			return
		}
	}

	// Enviar la respuesta con el usuario válido
	c.JSON(http.StatusOK, httputil.StandardResponse{
		Message: "Token válido",
		Data:    userInfo,
	})
}

// GetUserInfo obtiene la información del usuario desde Firebase Authentication usando el uid.
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

// needToUpdateClaims determina si es necesario actualizar los claims del token.
func needToUpdateClaims(claims map[string]interface{}, userInfo map[string]interface{}) bool {
	return claims["name"] != userInfo["displayName"] || claims["picture"] != userInfo["photoURL"]
}

// updateClaims actualiza los claims del token de usuario.
func updateClaims(authClient *auth.Client, uid string, userInfo map[string]interface{}) error {
	claims := map[string]interface{}{
		"name":    userInfo["displayName"],
		"picture": userInfo["photoURL"],
	}
	return authClient.SetCustomUserClaims(context.Background(), uid, claims)
}

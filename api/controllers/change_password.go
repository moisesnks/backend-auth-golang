package controllers

import (
	"backend/api/httputil"
	"backend/api/utils"
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// ChangePassword actualiza la contraseña de un usuario dado su UID
//
// Este endpoint permite actualizar la contraseña de un usuario autenticado utilizando Firebase Authentication.
//
// @Summary Actualizar contraseña
// @Description Actualiza la contraseña del usuario autenticado usando Firebase Authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Bearer Token"
// @Param body body ChangePasswordRequest true "Nueva contraseña del usuario"
// @Success 200 {object} httputil.StandardResponse "Contraseña actualizada correctamente"
// @Failure 400 {object} httputil.ErrorResponse "Datos de solicitud inválidos"
// @Failure 401 {object} httputil.ErrorResponse "Token de autenticación no proporcionado o inválido"
// @Failure 403 {object} httputil.ErrorResponse "Token de autenticación expirado"
// @Failure 500 {object} httputil.ErrorResponse "Error interno del servidor"
// @Router /change-password [post]
func ChangePassword(c *gin.Context, authClient *auth.Client) {
	// Obtener el cliente del middleware JWT
	claims, _ := c.Get("claims")
	if claims == nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "No se pudieron obtener los claims del token"})
		return
	}

	// Obtener el UID del usuario
	uid := claims.(*utils.Claims).UID

	// Verificar la fecha de expiración del token
	if claims.(*utils.Claims).IsTokenExpired() {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Token expirado"})
		return
	}

	// Obtener la contraseña del body
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: err.Error()})
		return
	}

	// Actualizar la contraseña usando el cliente de autenticación Firebase
	_, err := authClient.UpdateUser(context.Background(), uid, (&auth.UserToUpdate{}).Password(req.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Error al actualizar la contraseña"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, httputil.StandardResponse{
		Message: "Contraseña actualizada correctamente",
	})
}

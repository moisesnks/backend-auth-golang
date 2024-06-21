// auth/routes.go

package api

import (
	"github.com/gin-gonic/gin"

	"backend/api/controllers"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
)

// SetupRouter configura las rutas para el módulo de autenticación
func SetupRouter(r *gin.Engine, firestoreClient *firestore.Client, authClient *auth.Client) {
	authRoutes := r.Group("/")
	{
		authRoutes.POST("/login", controllers.LoginUser)
		authRoutes.POST("/register", func(c *gin.Context) {
			controllers.RegisterUser(c, firestoreClient)
		})
		authRoutes.POST("/verify", func(c *gin.Context) {
			controllers.VerifyCode(c, firestoreClient, authClient)
		})
	}
}

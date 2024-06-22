// auth/routes.go

package api

import (
	"github.com/gin-gonic/gin"

	"backend/api/controllers"
	"backend/api/middleware"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"firebase.google.com/go/auth"
)

// SetupRouter configura las rutas para el módulo de autenticación
func SetupRouter(r *gin.Engine, firestoreClient *firestore.Client, authClient *auth.Client, storageClient *storage.Client) {
	authRoutes := r.Group("/")
	{
		authRoutes.POST("/login", controllers.LoginUser)
		authRoutes.POST("/register", func(c *gin.Context) {
			controllers.RegisterUser(c, firestoreClient)
		})
		authRoutes.POST("/verify-code", func(c *gin.Context) {
			controllers.VerifyCode(c, firestoreClient, authClient)
		})
		authRoutes.POST("/resend-code", middleware.AuthMiddleware(authClient), func(c *gin.Context) {
			controllers.ResendCode(c, firestoreClient)
		})
		authRoutes.PATCH("/update", middleware.AuthMiddleware(authClient), func(c *gin.Context) {
			controllers.UpdateProfile(c, firestoreClient, authClient)
		})
		authRoutes.POST("/upload-photo", middleware.AuthMiddleware(authClient), func(c *gin.Context) {
			controllers.UploadPhoto(c, firestoreClient, storageClient, authClient)
		})
		authRoutes.POST("/forgot-password", func(c *gin.Context) {
			controllers.ForgotPassword(c, authClient, firestoreClient)
		})
		authRoutes.POST("/change-password", middleware.JWTMiddleware(), func(c *gin.Context) {
			controllers.ChangePassword(c, authClient)
		})
		authRoutes.GET("/validate-token", middleware.AuthMiddleware(authClient), func(c *gin.Context) {
			controllers.ValidateToken(c, authClient)
		})
	}
}

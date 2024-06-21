package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/api/option"

	"backend/api"
	_ "backend/docs"
)

var (
	firestoreClient *firestore.Client
	authClient      *auth.Client
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	// Obtener SERVICE_ACCOUNT_KEY de las variables de entorno que es un string without line breaks
	serviceAccountKey := os.Getenv("SERVICE_ACCOUNT_KEY")

	// Inicializar Firebase
	opt := option.WithCredentialsJSON([]byte(serviceAccountKey))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error inicializando app de Firebase: %v", err)
	}

	// Inicializar cliente de Firestore
	firestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error inicializando cliente de Firestore: %v", err)
	}

	// Inicializar cliente de Auth de Firebase
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error inicializando cliente de Auth de Firebase: %v", err)
	}

	// Cerrar cliente de Firestore al finalizar la aplicación
	defer firestoreClient.Close()

	// Configurar router Gin
	r := gin.New()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware JSON Logger
	r.Use(gin.Logger())

	// Middleware Recovery
	r.Use(gin.Recovery())

	// Configurar Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configurar rutas desde el paquete de autenticación (api)
	api.SetupRouter(r, firestoreClient, authClient)

	// Iniciar el servidor solo después de la inicialización completa
	port := "8080"
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

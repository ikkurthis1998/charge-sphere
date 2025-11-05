package main

import (
	"fmt"
	"log"

	"github.com/charge-sphere/ocpi-service/internal/api"
	"github.com/charge-sphere/ocpi-service/internal/api/handlers"
	"github.com/charge-sphere/ocpi-service/internal/api/middleware"
	"github.com/charge-sphere/ocpi-service/internal/config"
	"github.com/charge-sphere/ocpi-service/internal/database"
	"github.com/charge-sphere/ocpi-service/internal/domain/services"
	"github.com/charge-sphere/ocpi-service/internal/repository/mongodb"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MongoDB
	db, err := database.ConnectMongoDB(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	// Initialize repositories
	partnerRepo := mongodb.NewPartnerRepository(db)

	// Build hub URL
	hubURL := fmt.Sprintf("http://%s/ocpi/%s", cfg.GetServerAddr(), cfg.OCPI.Version)

	// Initialize services
	credentialsService := services.NewCredentialsService(
		partnerRepo,
		hubURL,
		cfg.OCPI.Version,
	)

	// Initialize handlers
	credentialsHandler := handlers.NewCredentialsHandler(credentialsService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(credentialsService)

	// Setup Gin router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router, credentialsHandler, authMiddleware)

	// Start server
	serverAddr := cfg.GetServerAddr()
	log.Printf("Starting OCPI service on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

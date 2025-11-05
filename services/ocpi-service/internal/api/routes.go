package api

import (
	"github.com/charge-sphere/ocpi-service/internal/api/handlers"
	"github.com/charge-sphere/ocpi-service/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	router *gin.Engine,
	credentialsHandler *handlers.CredentialsHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "ocpi-service",
		})
	})

	// OCPI API v2.3
	ocpi := router.Group("/ocpi/2.3")
	{
		// Credentials - no auth required for registration (POST)
		// but auth required for GET, PUT, DELETE
		credentials := ocpi.Group("/credentials")
		{
			// POST for registration - no auth required
			credentials.POST("", credentialsHandler.Register)

			// These require authentication
			authenticated := credentials.Group("")
			authenticated.Use(authMiddleware.Authenticate())
			{
				authenticated.GET("", credentialsHandler.Get)
				authenticated.PUT("", credentialsHandler.Update)
				authenticated.DELETE("", credentialsHandler.Delete)
			}
		}

		// Future modules will be added here
		// locations := ocpi.Group("/locations")
		// tokens := ocpi.Group("/tokens")
		// sessions := ocpi.Group("/sessions")
		// cdrs := ocpi.Group("/cdrs")
		// etc.
	}
}

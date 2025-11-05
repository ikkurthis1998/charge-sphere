package middleware

import (
	"net/http"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"github.com/charge-sphere/ocpi-service/internal/domain/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication for OCPI endpoints
type AuthMiddleware struct {
	credentialsService *services.CredentialsService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(credentialsService *services.CredentialsService) *AuthMiddleware {
	return &AuthMiddleware{
		credentialsService: credentialsService,
	}
}

// Authenticate validates the OCPI token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
				2001,
				"Missing authorization token",
				nil,
			))
			c.Abort()
			return
		}

		// Remove "Token " prefix if present
		if len(token) > 6 && token[:6] == "Token " {
			token = token[6:]
		}

		// Validate token
		partner, err := m.credentialsService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
				2001,
				"Invalid token: "+err.Error(),
				nil,
			))
			c.Abort()
			return
		}

		// Store partner info in context
		c.Set("partner", partner)
		c.Set("partner_id", partner.PartnerID)
		c.Set("partner_type", partner.Type)

		c.Next()
	}
}

// RequireCPO ensures the partner is a CPO
func (m *AuthMiddleware) RequireCPO() gin.HandlerFunc {
	return func(c *gin.Context) {
		partnerType, exists := c.Get("partner_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
				2001,
				"Partner type not found",
				nil,
			))
			c.Abort()
			return
		}

		if partnerType != models.PartnerTypeCPO {
			c.JSON(http.StatusForbidden, models.NewOCPIResponse(
				2001,
				"This endpoint requires CPO role",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireEMSP ensures the partner is an eMSP
func (m *AuthMiddleware) RequireEMSP() gin.HandlerFunc {
	return func(c *gin.Context) {
		partnerType, exists := c.Get("partner_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
				2001,
				"Partner type not found",
				nil,
			))
			c.Abort()
			return
		}

		if partnerType != models.PartnerTypeEMSP {
			c.JSON(http.StatusForbidden, models.NewOCPIResponse(
				2001,
				"This endpoint requires EMSP role",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

package handlers

import (
	"net/http"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"github.com/charge-sphere/ocpi-service/internal/domain/services"
	"github.com/gin-gonic/gin"
)

// CredentialsHandler handles credentials-related HTTP requests
type CredentialsHandler struct {
	credentialsService *services.CredentialsService
}

// NewCredentialsHandler creates a new credentials handler
func NewCredentialsHandler(credentialsService *services.CredentialsService) *CredentialsHandler {
	return &CredentialsHandler{
		credentialsService: credentialsService,
	}
}

// Register handles POST /ocpi/2.3/credentials
// Used by partners to register with the hub
func (h *CredentialsHandler) Register(c *gin.Context) {
	var req models.CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewOCPIResponse(
			1001,
			"Invalid request body: "+err.Error(),
			nil,
		))
		return
	}

	// Determine partner type from query param or header
	partnerType := models.PartnerType(c.Query("type"))
	if partnerType == "" {
		partnerType = models.PartnerType(c.GetHeader("X-Partner-Type"))
	}

	// Default to CPO if not specified
	if partnerType != models.PartnerTypeCPO && partnerType != models.PartnerTypeEMSP {
		partnerType = models.PartnerTypeCPO
	}

	credentials, err := h.credentialsService.RegisterPartner(c.Request.Context(), &req, partnerType)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewOCPIResponse(
			2001,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewOCPIResponse(
		1000,
		"Success",
		credentials,
	))
}

// Get handles GET /ocpi/2.3/credentials
// Returns hub credentials for the authenticated partner
func (h *CredentialsHandler) Get(c *gin.Context) {
	// Get token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
			2001,
			"Missing authorization token",
			nil,
		))
		return
	}

	// Remove "Token " prefix if present
	if len(token) > 6 && token[:6] == "Token " {
		token = token[6:]
	}

	credentials, err := h.credentialsService.GetCredentials(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
			2001,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewOCPIResponse(
		1000,
		"Success",
		credentials,
	))
}

// Update handles PUT /ocpi/2.3/credentials
// Updates partner credentials
func (h *CredentialsHandler) Update(c *gin.Context) {
	// Get token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
			2001,
			"Missing authorization token",
			nil,
		))
		return
	}

	// Remove "Token " prefix if present
	if len(token) > 6 && token[:6] == "Token " {
		token = token[6:]
	}

	var req models.CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewOCPIResponse(
			1001,
			"Invalid request body: "+err.Error(),
			nil,
		))
		return
	}

	credentials, err := h.credentialsService.UpdateCredentials(c.Request.Context(), token, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewOCPIResponse(
			2001,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewOCPIResponse(
		1000,
		"Success",
		credentials,
	))
}

// Delete handles DELETE /ocpi/2.3/credentials
// Removes partner registration
func (h *CredentialsHandler) Delete(c *gin.Context) {
	// Get token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.NewOCPIResponse(
			2001,
			"Missing authorization token",
			nil,
		))
		return
	}

	// Remove "Token " prefix if present
	if len(token) > 6 && token[:6] == "Token " {
		token = token[6:]
	}

	if err := h.credentialsService.DeleteCredentials(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, models.NewOCPIResponse(
			2001,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewOCPIResponse(
		1000,
		"Partner registration deleted successfully",
		nil,
	))
}

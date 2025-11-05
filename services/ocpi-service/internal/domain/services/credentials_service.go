package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
)

// PartnerRepository defines the interface for partner data operations
type PartnerRepository interface {
	Create(ctx context.Context, partner *models.Partner) error
	FindByPartnerID(ctx context.Context, partnerID string) (*models.Partner, error)
	FindByToken(ctx context.Context, token string) (*models.Partner, error)
	Update(ctx context.Context, partnerID string, partner *models.Partner) error
	Delete(ctx context.Context, partnerID string) error
	List(ctx context.Context, offset, limit int64) ([]*models.Partner, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, partnerID string, status models.PartnerStatus) error
}

// CredentialsService handles business logic for credentials
type CredentialsService struct {
	partnerRepo PartnerRepository
	hubURL      string
	version     string
}

// NewCredentialsService creates a new credentials service
func NewCredentialsService(partnerRepo PartnerRepository, hubURL string, version string) *CredentialsService {
	return &CredentialsService{
		partnerRepo: partnerRepo,
		hubURL:      hubURL,
		version:     version,
	}
}

// RegisterPartner registers a new partner (CPO or eMSP)
func (s *CredentialsService) RegisterPartner(ctx context.Context, req *models.CredentialsRequest, partnerType models.PartnerType) (*models.CredentialsResponse, error) {
	// Validate roles match partner type
	if err := s.validateRoles(req.Roles, partnerType); err != nil {
		return nil, err
	}

	// Generate a unique partner ID from the first role
	partnerID := fmt.Sprintf("%s-%s", req.Roles[0].CountryCode, req.Roles[0].PartyID)

	// Check if partner already exists
	existing, _ := s.partnerRepo.FindByPartnerID(ctx, partnerID)
	if existing != nil {
		return nil, fmt.Errorf("partner %s already registered", partnerID)
	}

	// Generate a new token for the hub
	hubToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create partner
	partner := &models.Partner{
		PartnerID: partnerID,
		Name:      s.extractBusinessName(req.Roles),
		Type:      partnerType,
		Credentials: models.Credentials{
			Token:   req.Token, // Partner's token for us to call them
			URL:     req.URL,
			Roles:   req.Roles,
			Version: s.version,
		},
		Status: models.PartnerStatusActive,
	}

	if err := s.partnerRepo.Create(ctx, partner); err != nil {
		return nil, err
	}

	// Return hub credentials for the partner to call us
	return &models.CredentialsResponse{
		Token: hubToken, // Our token for partner to call us
		URL:   s.hubURL,
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleHub,
				PartyID:     "HUB",
				CountryCode: "US",
				BusinessDetails: &models.BusinessDetails{
					Name: "ChargeSphere Hub",
				},
			},
		},
	}, nil
}

// GetCredentials returns the credentials for a partner to call the hub
func (s *CredentialsService) GetCredentials(ctx context.Context, token string) (*models.CredentialsResponse, error) {
	// Find partner by token
	partner, err := s.partnerRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	if partner.Status != models.PartnerStatusActive {
		return nil, fmt.Errorf("partner is not active")
	}

	// Return hub credentials
	return &models.CredentialsResponse{
		Token: token, // The partner's token to call us
		URL:   s.hubURL,
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleHub,
				PartyID:     "HUB",
				CountryCode: "US",
				BusinessDetails: &models.BusinessDetails{
					Name: "ChargeSphere Hub",
				},
			},
		},
	}, nil
}

// UpdateCredentials updates partner credentials
func (s *CredentialsService) UpdateCredentials(ctx context.Context, token string, req *models.CredentialsRequest) (*models.CredentialsResponse, error) {
	// Find partner by current token
	partner, err := s.partnerRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// Validate new roles match partner type
	if err := s.validateRoles(req.Roles, partner.Type); err != nil {
		return nil, err
	}

	// Update partner credentials
	partner.Credentials.Token = req.Token
	partner.Credentials.URL = req.URL
	partner.Credentials.Roles = req.Roles
	partner.Name = s.extractBusinessName(req.Roles)

	if err := s.partnerRepo.Update(ctx, partner.PartnerID, partner); err != nil {
		return nil, err
	}

	// Generate new hub token for updated credentials
	newHubToken, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.CredentialsResponse{
		Token: newHubToken,
		URL:   s.hubURL,
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleHub,
				PartyID:     "HUB",
				CountryCode: "US",
				BusinessDetails: &models.BusinessDetails{
					Name: "ChargeSphere Hub",
				},
			},
		},
	}, nil
}

// DeleteCredentials removes a partner registration
func (s *CredentialsService) DeleteCredentials(ctx context.Context, token string) error {
	// Find partner by token
	partner, err := s.partnerRepo.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// Delete partner
	if err := s.partnerRepo.Delete(ctx, partner.PartnerID); err != nil {
		return err
	}

	return nil
}

// ValidateToken validates a partner token and returns the partner
func (s *CredentialsService) ValidateToken(ctx context.Context, token string) (*models.Partner, error) {
	partner, err := s.partnerRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	if partner.Status != models.PartnerStatusActive {
		return nil, fmt.Errorf("partner is not active")
	}

	return partner, nil
}

// generateToken generates a secure random token
func (s *CredentialsService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// validateRoles validates that roles match the partner type
func (s *CredentialsService) validateRoles(roles []models.CredentialRole, partnerType models.PartnerType) error {
	if len(roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	for _, role := range roles {
		if partnerType == models.PartnerTypeCPO && role.Role != models.RoleCPO {
			return fmt.Errorf("CPO partner must have CPO role")
		}
		if partnerType == models.PartnerTypeEMSP && role.Role != models.RoleEMSP {
			return fmt.Errorf("eMSP partner must have EMSP role")
		}

		// Validate country code is 2 characters
		if len(role.CountryCode) != 2 {
			return fmt.Errorf("country code must be 2 characters")
		}

		// Validate party ID is 3 characters
		if len(role.PartyID) != 3 {
			return fmt.Errorf("party ID must be 3 characters")
		}
	}

	return nil
}

// extractBusinessName extracts business name from roles
func (s *CredentialsService) extractBusinessName(roles []models.CredentialRole) string {
	for _, role := range roles {
		if role.BusinessDetails != nil && role.BusinessDetails.Name != "" {
			return role.BusinessDetails.Name
		}
	}
	return fmt.Sprintf("%s-%s", roles[0].CountryCode, roles[0].PartyID)
}

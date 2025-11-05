package services

import (
	"context"
	"errors"
	"testing"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPartnerRepository is a mock implementation of PartnerRepository
type MockPartnerRepository struct {
	mock.Mock
}

func (m *MockPartnerRepository) Create(ctx context.Context, partner *models.Partner) error {
	args := m.Called(ctx, partner)
	return args.Error(0)
}

func (m *MockPartnerRepository) FindByPartnerID(ctx context.Context, partnerID string) (*models.Partner, error) {
	args := m.Called(ctx, partnerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Partner), args.Error(1)
}

func (m *MockPartnerRepository) FindByToken(ctx context.Context, token string) (*models.Partner, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Partner), args.Error(1)
}

func (m *MockPartnerRepository) Update(ctx context.Context, partnerID string, partner *models.Partner) error {
	args := m.Called(ctx, partnerID, partner)
	return args.Error(0)
}

func (m *MockPartnerRepository) Delete(ctx context.Context, partnerID string) error {
	args := m.Called(ctx, partnerID)
	return args.Error(0)
}

func (m *MockPartnerRepository) List(ctx context.Context, offset, limit int64) ([]*models.Partner, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Partner), args.Error(1)
}

func (m *MockPartnerRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPartnerRepository) UpdateStatus(ctx context.Context, partnerID string, status models.PartnerStatus) error {
	args := m.Called(ctx, partnerID, status)
	return args.Error(0)
}

func TestCredentialsService_RegisterPartner(t *testing.T) {
	mockRepo := new(MockPartnerRepository)
	service := NewCredentialsService(mockRepo, "http://localhost:8080/ocpi/2.3", "2.3")

	tests := []struct {
		name          string
		request       *models.CredentialsRequest
		partnerType   models.PartnerType
		setupMock     func()
		expectError   bool
		errorContains string
	}{
		{
			name: "successful CPO registration",
			request: &models.CredentialsRequest{
				Token: "partner_token_123",
				URL:   "https://partner.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "ABC",
						CountryCode: "DE",
						BusinessDetails: &models.BusinessDetails{
							Name: "Test CPO",
						},
					},
				},
			},
			partnerType: models.PartnerTypeCPO,
			setupMock: func() {
				mockRepo.On("FindByPartnerID", mock.Anything, "DE-ABC").Return(nil, errors.New("not found")).Once()
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Partner) bool {
					return p.PartnerID == "DE-ABC" && p.Type == models.PartnerTypeCPO
				})).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "successful eMSP registration",
			request: &models.CredentialsRequest{
				Token: "emsp_token_456",
				URL:   "https://emsp.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleEMSP,
						PartyID:     "XYZ",
						CountryCode: "US",
						BusinessDetails: &models.BusinessDetails{
							Name: "Test eMSP",
						},
					},
				},
			},
			partnerType: models.PartnerTypeEMSP,
			setupMock: func() {
				mockRepo.On("FindByPartnerID", mock.Anything, "US-XYZ").Return(nil, errors.New("not found")).Once()
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Partner) bool {
					return p.PartnerID == "US-XYZ" && p.Type == models.PartnerTypeEMSP
				})).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "partner already exists",
			request: &models.CredentialsRequest{
				Token: "partner_token_789",
				URL:   "https://existing.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "EXS",
						CountryCode: "FR",
					},
				},
			},
			partnerType: models.PartnerTypeCPO,
			setupMock: func() {
				existingPartner := &models.Partner{
					PartnerID: "FR-EXS",
					Status:    models.PartnerStatusActive,
				}
				mockRepo.On("FindByPartnerID", mock.Anything, "FR-EXS").Return(existingPartner, nil).Once()
			},
			expectError:   true,
			errorContains: "already registered",
		},
		{
			name: "invalid role for CPO",
			request: &models.CredentialsRequest{
				Token: "invalid_token",
				URL:   "https://invalid.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleEMSP, // Wrong role for CPO
						PartyID:     "INV",
						CountryCode: "UK",
					},
				},
			},
			partnerType:   models.PartnerTypeCPO,
			setupMock:     func() {},
			expectError:   true,
			errorContains: "must have CPO role",
		},
		{
			name: "invalid country code",
			request: &models.CredentialsRequest{
				Token: "invalid_token",
				URL:   "https://invalid.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "INV",
						CountryCode: "USA", // Should be 2 characters
					},
				},
			},
			partnerType:   models.PartnerTypeCPO,
			setupMock:     func() {},
			expectError:   true,
			errorContains: "country code must be 2 characters",
		},
		{
			name: "invalid party ID",
			request: &models.CredentialsRequest{
				Token: "invalid_token",
				URL:   "https://invalid.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "AB", // Should be 3 characters
						CountryCode: "US",
					},
				},
			},
			partnerType:   models.PartnerTypeCPO,
			setupMock:     func() {},
			expectError:   true,
			errorContains: "party ID must be 3 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			response, err := service.RegisterPartner(context.Background(), tt.request, tt.partnerType)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, "http://localhost:8080/ocpi/2.3", response.URL)
				assert.Len(t, response.Roles, 1)
				assert.Equal(t, models.RoleHub, response.Roles[0].Role)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCredentialsService_GetCredentials(t *testing.T) {
	mockRepo := new(MockPartnerRepository)
	service := NewCredentialsService(mockRepo, "http://localhost:8080/ocpi/2.3", "2.3")

	tests := []struct {
		name          string
		token         string
		setupMock     func()
		expectError   bool
		errorContains string
	}{
		{
			name:  "successful get credentials",
			token: "valid_token_123",
			setupMock: func() {
				partner := &models.Partner{
					PartnerID: "DE-ABC",
					Status:    models.PartnerStatusActive,
					Credentials: models.Credentials{
						Token: "valid_token_123",
					},
				}
				mockRepo.On("FindByToken", mock.Anything, "valid_token_123").Return(partner, nil).Once()
			},
			expectError: false,
		},
		{
			name:  "token not found",
			token: "invalid_token",
			setupMock: func() {
				mockRepo.On("FindByToken", mock.Anything, "invalid_token").Return(nil, errors.New("not found")).Once()
			},
			expectError:   true,
			errorContains: "invalid token",
		},
		{
			name:  "inactive partner",
			token: "inactive_token",
			setupMock: func() {
				partner := &models.Partner{
					PartnerID: "FR-XYZ",
					Status:    models.PartnerStatusInactive,
					Credentials: models.Credentials{
						Token: "inactive_token",
					},
				}
				mockRepo.On("FindByToken", mock.Anything, "inactive_token").Return(partner, nil).Once()
			},
			expectError:   true,
			errorContains: "not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			response, err := service.GetCredentials(context.Background(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.token, response.Token)
				assert.Equal(t, "http://localhost:8080/ocpi/2.3", response.URL)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCredentialsService_ValidateToken(t *testing.T) {
	mockRepo := new(MockPartnerRepository)
	service := NewCredentialsService(mockRepo, "http://localhost:8080/ocpi/2.3", "2.3")

	tests := []struct {
		name          string
		token         string
		setupMock     func()
		expectError   bool
		errorContains string
	}{
		{
			name:  "valid token",
			token: "valid_token",
			setupMock: func() {
				partner := &models.Partner{
					PartnerID: "DE-ABC",
					Status:    models.PartnerStatusActive,
				}
				mockRepo.On("FindByToken", mock.Anything, "valid_token").Return(partner, nil).Once()
			},
			expectError: false,
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			setupMock: func() {
				mockRepo.On("FindByToken", mock.Anything, "invalid_token").Return(nil, errors.New("not found")).Once()
			},
			expectError:   true,
			errorContains: "invalid token",
		},
		{
			name:  "suspended partner",
			token: "suspended_token",
			setupMock: func() {
				partner := &models.Partner{
					PartnerID: "US-XYZ",
					Status:    models.PartnerStatusSuspended,
				}
				mockRepo.On("FindByToken", mock.Anything, "suspended_token").Return(partner, nil).Once()
			},
			expectError:   true,
			errorContains: "not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			partner, err := service.ValidateToken(context.Background(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, partner)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, partner)
				assert.Equal(t, models.PartnerStatusActive, partner.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCredentialsService_DeleteCredentials(t *testing.T) {
	mockRepo := new(MockPartnerRepository)
	service := NewCredentialsService(mockRepo, "http://localhost:8080/ocpi/2.3", "2.3")

	t.Run("successful deletion", func(t *testing.T) {
		partner := &models.Partner{
			PartnerID: "DE-ABC",
			Status:    models.PartnerStatusActive,
		}
		mockRepo.On("FindByToken", mock.Anything, "valid_token").Return(partner, nil).Once()
		mockRepo.On("Delete", mock.Anything, "DE-ABC").Return(nil).Once()

		err := service.DeleteCredentials(context.Background(), "valid_token")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("token not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("FindByToken", mock.Anything, "invalid_token").Return(nil, errors.New("not found")).Once()

		err := service.DeleteCredentials(context.Background(), "invalid_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		mockRepo.AssertExpectations(t)
	})
}

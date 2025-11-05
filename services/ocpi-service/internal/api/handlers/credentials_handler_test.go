package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCredentialsService is a mock implementation of CredentialsService
type MockCredentialsService struct {
	mock.Mock
}

func (m *MockCredentialsService) RegisterPartner(ctx mock.Anything, req *models.CredentialsRequest, partnerType models.PartnerType) (*models.CredentialsResponse, error) {
	args := m.Called(ctx, req, partnerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CredentialsResponse), args.Error(1)
}

func (m *MockCredentialsService) GetCredentials(ctx mock.Anything, token string) (*models.CredentialsResponse, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CredentialsResponse), args.Error(1)
}

func (m *MockCredentialsService) UpdateCredentials(ctx mock.Anything, token string, req *models.CredentialsRequest) (*models.CredentialsResponse, error) {
	args := m.Called(ctx, token, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CredentialsResponse), args.Error(1)
}

func (m *MockCredentialsService) DeleteCredentials(ctx mock.Anything, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockCredentialsService) ValidateToken(ctx mock.Anything, token string) (*models.Partner, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Partner), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCredentialsHandler_Register(t *testing.T) {
	mockService := new(MockCredentialsService)
	handler := NewCredentialsHandler(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		queryParams    map[string]string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "successful CPO registration",
			requestBody: models.CredentialsRequest{
				Token: "partner_token",
				URL:   "https://partner.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "ABC",
						CountryCode: "DE",
					},
				},
			},
			queryParams: map[string]string{
				"type": "CPO",
			},
			setupMock: func() {
				response := &models.CredentialsResponse{
					Token: "hub_token_123",
					URL:   "http://localhost/ocpi/2.3",
					Roles: []models.CredentialRole{
						{
							Role:        models.RoleHub,
							PartyID:     "HUB",
							CountryCode: "US",
						},
					},
				}
				mockService.On("RegisterPartner", mock.Anything, mock.Anything, models.PartnerTypeCPO).Return(response, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1000), response["status_code"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "hub_token_123", data["token"])
			},
		},
		{
			name: "successful eMSP registration",
			requestBody: models.CredentialsRequest{
				Token: "emsp_token",
				URL:   "https://emsp.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleEMSP,
						PartyID:     "XYZ",
						CountryCode: "US",
					},
				},
			},
			queryParams: map[string]string{
				"type": "EMSP",
			},
			setupMock: func() {
				response := &models.CredentialsResponse{
					Token: "hub_token_456",
					URL:   "http://localhost/ocpi/2.3",
					Roles: []models.CredentialRole{
						{
							Role:        models.RoleHub,
							PartyID:     "HUB",
							CountryCode: "US",
						},
					},
				}
				mockService.On("RegisterPartner", mock.Anything, mock.Anything, models.PartnerTypeEMSP).Return(response, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1000), response["status_code"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "hub_token_456", data["token"])
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			queryParams:    map[string]string{},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1001), response["status_code"])
			},
		},
		{
			name: "registration error",
			requestBody: models.CredentialsRequest{
				Token: "partner_token",
				URL:   "https://partner.com/ocpi",
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     "ABC",
						CountryCode: "DE",
					},
				},
			},
			queryParams: map[string]string{
				"type": "CPO",
			},
			setupMock: func() {
				mockService.On("RegisterPartner", mock.Anything, mock.Anything, models.PartnerTypeCPO).Return(nil, errors.New("partner already exists")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(2001), response["status_code"])
				assert.Contains(t, response["status_message"], "already exists")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			router := setupTestRouter()
			router.POST("/credentials", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/credentials", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Add query params
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCredentialsHandler_Get(t *testing.T) {
	mockService := new(MockCredentialsService)
	handler := NewCredentialsHandler(mockService)

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, response map[string]interface{})
	}{
		{
			name:       "successful get with Token prefix",
			authHeader: "Token valid_token_123",
			setupMock: func() {
				response := &models.CredentialsResponse{
					Token: "valid_token_123",
					URL:   "http://localhost/ocpi/2.3",
					Roles: []models.CredentialRole{
						{
							Role:        models.RoleHub,
							PartyID:     "HUB",
							CountryCode: "US",
						},
					},
				}
				mockService.On("GetCredentials", mock.Anything, "valid_token_123").Return(response, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1000), response["status_code"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "valid_token_123", data["token"])
			},
		},
		{
			name:       "successful get without Token prefix",
			authHeader: "raw_token_456",
			setupMock: func() {
				response := &models.CredentialsResponse{
					Token: "raw_token_456",
					URL:   "http://localhost/ocpi/2.3",
					Roles: []models.CredentialRole{
						{
							Role:        models.RoleHub,
							PartyID:     "HUB",
							CountryCode: "US",
						},
					},
				}
				mockService.On("GetCredentials", mock.Anything, "raw_token_456").Return(response, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1000), response["status_code"])
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(2001), response["status_code"])
				assert.Contains(t, response["status_message"], "Missing authorization")
			},
		},
		{
			name:       "invalid token",
			authHeader: "Token invalid_token",
			setupMock: func() {
				mockService.On("GetCredentials", mock.Anything, "invalid_token").Return(nil, errors.New("invalid token")).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(2001), response["status_code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			router := setupTestRouter()
			router.GET("/credentials", handler.Get)

			req, _ := http.NewRequest(http.MethodGet, "/credentials", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCredentialsHandler_Delete(t *testing.T) {
	mockService := new(MockCredentialsService)
	handler := NewCredentialsHandler(mockService)

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, response map[string]interface{})
	}{
		{
			name:       "successful deletion",
			authHeader: "Token valid_token",
			setupMock: func() {
				mockService.On("DeleteCredentials", mock.Anything, "valid_token").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(1000), response["status_code"])
				assert.Contains(t, response["status_message"], "deleted successfully")
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(2001), response["status_code"])
			},
		},
		{
			name:       "deletion error",
			authHeader: "Token invalid_token",
			setupMock: func() {
				mockService.On("DeleteCredentials", mock.Anything, "invalid_token").Return(errors.New("partner not found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, float64(2001), response["status_code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			router := setupTestRouter()
			router.DELETE("/credentials", handler.Delete)

			req, _ := http.NewRequest(http.MethodDelete, "/credentials", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

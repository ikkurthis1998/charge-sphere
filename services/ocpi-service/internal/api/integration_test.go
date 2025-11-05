package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/charge-sphere/ocpi-service/internal/api"
	"github.com/charge-sphere/ocpi-service/internal/api/handlers"
	"github.com/charge-sphere/ocpi-service/internal/api/middleware"
	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"github.com/charge-sphere/ocpi-service/internal/domain/services"
	"github.com/charge-sphere/ocpi-service/internal/repository/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IntegrationTestSuite is the test suite for integration tests
type IntegrationTestSuite struct {
	suite.Suite
	mongoClient *mongo.Client
	db          *mongo.Database
	router      *gin.Engine
	cleanup     func()
}

// SetupSuite runs once before all tests
func (suite *IntegrationTestSuite) SetupSuite() {
	// Connect to test MongoDB
	// You can use a test container or local MongoDB
	// For this example, we'll use the same MongoDB but a different database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := "mongodb://admin:admin123@localhost:27017"
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	require.NoError(suite.T(), err, "Failed to connect to MongoDB")

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	require.NoError(suite.T(), err, "Failed to ping MongoDB")

	suite.mongoClient = client
	suite.db = client.Database("chargesphere_test")
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	// Clean up collections before each test
	ctx := context.Background()
	suite.db.Collection("partners").Drop(ctx)

	// Initialize repositories
	partnerRepo := mongodb.NewPartnerRepository(suite.db)

	// Initialize services
	credentialsService := services.NewCredentialsService(
		partnerRepo,
		"http://localhost:8080/ocpi/2.3",
		"2.3",
	)

	// Initialize handlers
	credentialsHandler := handlers.NewCredentialsHandler(credentialsService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(credentialsService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupRoutes(router, credentialsHandler, authMiddleware)

	suite.router = router
}

// TearDownSuite runs once after all tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.mongoClient != nil {
		ctx := context.Background()
		suite.db.Drop(ctx)
		suite.mongoClient.Disconnect(ctx)
	}
}

// TestIntegrationSuite runs the test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Helper function to make HTTP requests
func (suite *IntegrationTestSuite) makeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

// Test: Complete partner registration flow
func (suite *IntegrationTestSuite) TestCompleteRegistrationFlow() {
	t := suite.T()

	// Step 1: Register a CPO
	registerReq := models.CredentialsRequest{
		Token: "cpo_token_123",
		URL:   "https://cpo-company.com/ocpi/2.3",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "ABC",
				CountryCode: "DE",
				BusinessDetails: &models.BusinessDetails{
					Name: "ABC Charging Network",
				},
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var registerResponse models.OCPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	require.NoError(t, err)

	assert.Equal(t, 1000, registerResponse.StatusCode)
	assert.NotNil(t, registerResponse.Data)

	// Extract hub token from response
	responseData := registerResponse.Data.(map[string]interface{})
	hubToken := responseData["token"].(string)
	assert.NotEmpty(t, hubToken)

	// Step 2: Use hub token to get credentials
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + hubToken,
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse models.OCPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	require.NoError(t, err)

	assert.Equal(t, 1000, getResponse.StatusCode)
	getData := getResponse.Data.(map[string]interface{})
	assert.Equal(t, hubToken, getData["token"].(string))

	// Step 3: Update credentials
	updateReq := models.CredentialsRequest{
		Token: "cpo_token_updated_456",
		URL:   "https://cpo-company-new.com/ocpi/2.3",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "ABC",
				CountryCode: "DE",
				BusinessDetails: &models.BusinessDetails{
					Name: "ABC Charging Network Updated",
				},
			},
		},
	}

	w = suite.makeRequest("PUT", "/ocpi/2.3/credentials", updateReq, map[string]string{
		"Authorization": "Token " + hubToken,
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse models.OCPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	require.NoError(t, err)

	assert.Equal(t, 1000, updateResponse.StatusCode)
	updateData := updateResponse.Data.(map[string]interface{})
	newHubToken := updateData["token"].(string)
	assert.NotEmpty(t, newHubToken)
	assert.NotEqual(t, hubToken, newHubToken) // Token should be regenerated

	// Step 4: Old token should still work (tokens are not invalidated immediately in this implementation)
	// In production, you might want to invalidate old tokens

	// Step 5: Delete registration using new token
	w = suite.makeRequest("DELETE", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + newHubToken,
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var deleteResponse models.OCPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &deleteResponse)
	require.NoError(t, err)

	assert.Equal(t, 1000, deleteResponse.StatusCode)
	assert.Contains(t, deleteResponse.StatusMessage, "deleted successfully")

	// Step 6: Try to use deleted token - should fail
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + newHubToken,
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: Register multiple partners (CPO and eMSP)
func (suite *IntegrationTestSuite) TestRegisterMultiplePartners() {
	t := suite.T()

	// Register CPO
	cpoReq := models.CredentialsRequest{
		Token: "cpo_token",
		URL:   "https://cpo.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "CPO",
				CountryCode: "DE",
				BusinessDetails: &models.BusinessDetails{
					Name: "CPO Company",
				},
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", cpoReq, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var cpoResponse models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &cpoResponse)
	cpoData := cpoResponse.Data.(map[string]interface{})
	cpoToken := cpoData["token"].(string)

	// Register eMSP
	emspReq := models.CredentialsRequest{
		Token: "emsp_token",
		URL:   "https://emsp.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleEMSP,
				PartyID:     "EMS",
				CountryCode: "US",
				BusinessDetails: &models.BusinessDetails{
					Name: "eMSP Company",
				},
			},
		},
	}

	w = suite.makeRequest("POST", "/ocpi/2.3/credentials?type=EMSP", emspReq, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var emspResponse models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &emspResponse)
	emspData := emspResponse.Data.(map[string]interface{})
	emspToken := emspData["token"].(string)

	// Both tokens should be different
	assert.NotEqual(t, cpoToken, emspToken)

	// Both should be able to get their credentials
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + cpoToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + emspToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test: Duplicate partner registration should fail
func (suite *IntegrationTestSuite) TestDuplicatePartnerRegistration() {
	t := suite.T()

	registerReq := models.CredentialsRequest{
		Token: "token_123",
		URL:   "https://partner.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "DUP",
				CountryCode: "FR",
			},
		},
	}

	// First registration should succeed
	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// Second registration with same partner ID should fail
	w = suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response.StatusMessage, "already registered")
}

// Test: Invalid token authentication
func (suite *IntegrationTestSuite) TestInvalidTokenAuthentication() {
	t := suite.T()

	// Try to get credentials with invalid token
	w := suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token invalid_token_xyz",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 2001, response.StatusCode)
	assert.Contains(t, response.StatusMessage, "Invalid token")
}

// Test: Missing authorization header
func (suite *IntegrationTestSuite) TestMissingAuthorizationHeader() {
	t := suite.T()

	w := suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 2001, response.StatusCode)
	assert.Contains(t, response.StatusMessage, "Missing authorization")
}

// Test: Invalid request body
func (suite *IntegrationTestSuite) TestInvalidRequestBody() {
	t := suite.T()

	invalidReq := map[string]interface{}{
		"invalid": "data",
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials", invalidReq, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1001, response.StatusCode)
}

// Test: Invalid role for partner type
func (suite *IntegrationTestSuite) TestInvalidRoleForPartnerType() {
	t := suite.T()

	// Try to register CPO with eMSP role
	registerReq := models.CredentialsRequest{
		Token: "token_123",
		URL:   "https://partner.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleEMSP, // Wrong role for CPO type
				PartyID:     "WRG",
				CountryCode: "UK",
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response.StatusMessage, "must have CPO role")
}

// Test: Invalid country code (not 2 characters)
func (suite *IntegrationTestSuite) TestInvalidCountryCode() {
	t := suite.T()

	registerReq := models.CredentialsRequest{
		Token: "token_123",
		URL:   "https://partner.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "ABC",
				CountryCode: "USA", // Should be 2 characters
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response.StatusMessage, "country code must be 2 characters")
}

// Test: Invalid party ID (not 3 characters)
func (suite *IntegrationTestSuite) TestInvalidPartyID() {
	t := suite.T()

	registerReq := models.CredentialsRequest{
		Token: "token_123",
		URL:   "https://partner.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "AB", // Should be 3 characters
				CountryCode: "US",
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response.StatusMessage, "party ID must be 3 characters")
}

// Test: Health check endpoint
func (suite *IntegrationTestSuite) TestHealthCheck() {
	t := suite.T()

	w := suite.makeRequest("GET", "/health", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "ocpi-service", response["service"])
}

// Test: Data persistence across requests
func (suite *IntegrationTestSuite) TestDataPersistence() {
	t := suite.T()

	// Register a partner
	registerReq := models.CredentialsRequest{
		Token: "persist_token",
		URL:   "https://persist.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "PER",
				CountryCode: "IT",
				BusinessDetails: &models.BusinessDetails{
					Name:    "Persistent Partner",
					Website: "https://persist.com",
				},
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)
	require.Equal(t, http.StatusOK, w.Code)

	var registerResponse models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	responseData := registerResponse.Data.(map[string]interface{})
	hubToken := responseData["token"].(string)

	// Verify data was persisted by querying MongoDB directly
	ctx := context.Background()
	var partner models.Partner
	err := suite.db.Collection("partners").FindOne(ctx, map[string]interface{}{
		"partner_id": "IT-PER",
	}).Decode(&partner)

	require.NoError(t, err)
	assert.Equal(t, "IT-PER", partner.PartnerID)
	assert.Equal(t, "Persistent Partner", partner.Name)
	assert.Equal(t, models.PartnerTypeCPO, partner.Type)
	assert.Equal(t, models.PartnerStatusActive, partner.Status)
	assert.Equal(t, "persist_token", partner.Credentials.Token)
	assert.Len(t, partner.Credentials.Roles, 1)
	assert.Equal(t, models.RoleCPO, partner.Credentials.Roles[0].Role)

	// Get credentials and verify response
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + hubToken,
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test: Concurrent partner registrations
func (suite *IntegrationTestSuite) TestConcurrentRegistrations() {
	t := suite.T()

	numPartners := 10
	done := make(chan bool, numPartners)
	errors := make(chan error, numPartners)

	for i := 0; i < numPartners; i++ {
		go func(index int) {
			registerReq := models.CredentialsRequest{
				Token: fmt.Sprintf("token_%d", index),
				URL:   fmt.Sprintf("https://partner%d.com/ocpi", index),
				Roles: []models.CredentialRole{
					{
						Role:        models.RoleCPO,
						PartyID:     fmt.Sprintf("P%02d", index),
						CountryCode: "DE",
					},
				},
			}

			w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf("registration %d failed with status %d", index, w.Code)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numPartners; i++ {
		<-done
	}
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Verify all partners were created
	ctx := context.Background()
	count, err := suite.db.Collection("partners").CountDocuments(ctx, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, int64(numPartners), count)
}

// Test: Token with and without "Token " prefix
func (suite *IntegrationTestSuite) TestTokenPrefixHandling() {
	t := suite.T()

	// Register partner
	registerReq := models.CredentialsRequest{
		Token: "prefix_test_token",
		URL:   "https://prefix.com/ocpi",
		Roles: []models.CredentialRole{
			{
				Role:        models.RoleCPO,
				PartyID:     "PRE",
				CountryCode: "ES",
			},
		},
	}

	w := suite.makeRequest("POST", "/ocpi/2.3/credentials?type=CPO", registerReq, nil)
	require.Equal(t, http.StatusOK, w.Code)

	var response models.OCPIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	responseData := response.Data.(map[string]interface{})
	hubToken := responseData["token"].(string)

	// Test with "Token " prefix
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": "Token " + hubToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	// Test without "Token " prefix
	w = suite.makeRequest("GET", "/ocpi/2.3/credentials", nil, map[string]string{
		"Authorization": hubToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

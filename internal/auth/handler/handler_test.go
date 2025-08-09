package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chinese-bridge-game/internal/auth/dto"
	"chinese-bridge-game/internal/auth/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GoogleOAuthLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, tokenString string) (*dto.JWTClaims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JWTClaims), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) GetGoogleOAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func setupTestRouter(authService service.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add trace ID middleware for testing
	router.Use(func(c *gin.Context) {
		c.Set("trace_id", "test-trace-id")
		c.Next()
	})
	
	handler := NewAuthHandler(authService)
	api := router.Group("/api/v1")
	
	// Add health endpoints
	api.GET("/health", handler.HealthCheck)
	api.GET("/ready", handler.ReadyCheck)
	
	handler.RegisterRoutes(api)
	
	return router
}

func TestAuthHandler_GetGoogleOAuthURL(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Setup expectations
	expectedURL := "https://accounts.google.com/oauth2/auth?client_id=test&redirect_uri=test&response_type=code&scope=email+profile&state=test-state"
	mockService.On("GetGoogleOAuthURL", "test-state").Return(expectedURL)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/auth/google/url?state=test-state", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, response["url"])
	assert.Equal(t, "test-state", response["state"])

	mockService.AssertExpectations(t)
}

func TestAuthHandler_GoogleOAuthCallback_Success(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Test data
	request := dto.GoogleOAuthRequest{
		Code:  "test-auth-code",
		State: "test-state",
	}

	expectedResponse := &dto.AuthResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: dto.UserInfo{
			ID:       "test-user-id",
			GoogleID: "test-google-id",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "https://example.com/avatar.jpg",
		},
	}

	// Setup expectations
	mockService.On("GoogleOAuthLogin", mock.Anything, "test-auth-code").Return(expectedResponse, nil)

	// Create request
	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/auth/google", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
	assert.Equal(t, expectedResponse.RefreshToken, response.RefreshToken)
	assert.Equal(t, expectedResponse.User.Email, response.User.Email)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_GoogleOAuthCallback_InvalidRequest(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Create invalid request (missing required code field)
	invalidRequest := map[string]string{
		"state": "test-state",
	}

	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("POST", "/api/v1/auth/google", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_ERROR", response.Code)
	assert.Contains(t, response.Message, "Invalid request body")
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Test data
	request := dto.RefreshTokenRequest{
		RefreshToken: "test-refresh-token",
	}

	expectedResponse := &dto.TokenResponse{
		AccessToken: "new-access-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	// Setup expectations
	mockService.On("RefreshToken", mock.Anything, "test-refresh-token").Return(expectedResponse, nil)

	// Create request
	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
	assert.Equal(t, expectedResponse.TokenType, response.TokenType)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Setup JWT validation mock
	mockService.On("ValidateToken", mock.Anything, "valid-token").Return(&dto.JWTClaims{
		UserID: "test-user-id",
		Email:  "test@example.com",
		Name:   "Test User",
	}, nil)

	// Setup logout mock
	mockService.On("Logout", mock.Anything, "test-user-id").Return(nil)

	// Create request with valid JWT token
	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully logged out", response.Message)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Unauthorized(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Create request without authorization header
	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "AUTHENTICATION_ERROR", response.Code)
}

func TestAuthHandler_HealthCheck(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "auth-service", response["service"])
}

func TestAuthHandler_ReadyCheck(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/ready", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ready", response["status"])
	assert.Equal(t, "auth-service", response["service"])
}

// Test rate limiting
func TestAuthHandler_RateLimit(t *testing.T) {
	// Setup
	mockService := new(MockAuthService)
	router := setupTestRouter(mockService)

	// Setup mock expectations for all requests
	expectedURL := "https://accounts.google.com/oauth2/auth?test=true"
	mockService.On("GetGoogleOAuthURL", mock.AnythingOfType("string")).Return(expectedURL)

	// Make multiple requests quickly to trigger rate limit
	for i := 0; i < 15; i++ { // Exceed the rate limit of 10 requests per burst
		req, _ := http.NewRequest("GET", "/api/v1/auth/google/url", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if i < 10 {
			// First 10 requests should succeed (within burst limit)
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusTooManyRequests)
		} else {
			// Subsequent requests should be rate limited
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}
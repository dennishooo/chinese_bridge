package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"chinese-bridge-game/internal/auth/dto"
	"chinese-bridge-game/internal/common/config"
	"chinese-bridge-game/internal/common/database"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *database.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*database.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockAuthRepository) UpdateUser(ctx context.Context, user *database.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) CreateSession(ctx context.Context, session *database.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepository) GetSessionByToken(ctx context.Context, token string) (*database.Session, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.Session), args.Error(1)
}

func (m *MockAuthRepository) DeleteSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteUserSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockRedisClient is a mock implementation of RedisClient interface
type MockRedisClient struct {
	mock.Mock
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	// Handle both string and []byte values
	switch v := value.(type) {
	case string:
		m.data[key] = v
	case []byte:
		m.data[key] = string(v)
	default:
		m.data[key] = fmt.Sprintf("%v", v)
	}
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if value, exists := m.data[key]; exists {
		cmd.SetVal(value)
	} else {
		cmd.SetErr(redis.Nil)
	}
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	for _, key := range keys {
		delete(m.data, key)
	}
	cmd := redis.NewIntCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	args := m.Called(ctx, pattern)
	var keys []string
	for key := range m.data {
		keys = append(keys, key)
	}
	cmd := redis.NewStringSliceCmd(ctx)
	cmd.SetVal(keys)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(mockRepo, mockRedis, cfg).(*authService)

	// Test data
	user := &database.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate a valid token
	token, err := service.generateAccessToken(user)
	assert.NoError(t, err)

	// Test valid token
	claims, err := service.ValidateToken(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Name, claims.Name)

	// Test invalid token
	_, err = service.ValidateToken(context.Background(), "invalid-token")
	assert.Error(t, err)

	// Test expired token
	expiredClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString([]byte(cfg.JWTSecret))
	assert.NoError(t, err)

	_, err = service.ValidateToken(context.Background(), expiredTokenString)
	assert.Error(t, err)
}

func TestAuthService_GenerateAccessToken(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(mockRepo, mockRedis, cfg).(*authService)

	// Test data
	user := &database.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate token
	tokenString, err := service.generateAccessToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Verify token can be parsed
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, user.ID, claims["user_id"])
	assert.Equal(t, user.Email, claims["email"])
	assert.Equal(t, user.Name, claims["name"])
}

func TestAuthService_GenerateRefreshToken(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(mockRepo, mockRedis, cfg).(*authService)

	// Generate refresh token
	token1, err := service.generateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Generate another refresh token
	token2, err := service.generateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Tokens should be different
	assert.NotEqual(t, token1, token2)

	// Tokens should be base64 URL encoded (can contain = padding)
	assert.Regexp(t, `^[A-Za-z0-9_-]+=*$`, token1)
	assert.Regexp(t, `^[A-Za-z0-9_-]+=*$`, token2)
}

func TestAuthService_Logout(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(mockRepo, mockRedis, cfg)

	userID := "test-user-id"

	// Setup expectations
	mockRepo.On("DeleteUserSessions", mock.Anything, userID).Return(nil)
	mockRedis.On("Keys", mock.Anything, "session:*").Return(nil)

	// Test logout
	err := service.Logout(context.Background(), userID)
	assert.NoError(t, err)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetGoogleOAuthURL(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		GoogleOAuth: config.GoogleOAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/auth/google/callback",
		},
	}

	service := NewAuthService(mockRepo, mockRedis, cfg)

	// Test URL generation
	state := "test-state"
	url := service.GetGoogleOAuthURL(state)

	assert.NotEmpty(t, url)
	assert.Contains(t, url, "accounts.google.com")
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fauth%2Fgoogle%2Fcallback")
}

// Integration test helper functions
func setupTestService() (*authService, *MockAuthRepository, *MockRedisClient) {
	mockRepo := new(MockAuthRepository)
	mockRedis := NewMockRedisClient()
	cfg := &config.Config{
		JWTSecret: "test-secret",
		GoogleOAuth: config.GoogleOAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/auth/google/callback",
		},
	}

	service := NewAuthService(mockRepo, mockRedis, cfg).(*authService)
	return service, mockRepo, mockRedis
}

func TestAuthService_StoreAndGetSession(t *testing.T) {
	service, _, mockRedis := setupTestService()

	// Test data
	refreshToken := "test-refresh-token"
	sessionInfo := &dto.SessionInfo{
		UserID:       "test-user-id",
		Email:        "test@example.com",
		Name:         "Test User",
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Setup expectations
	mockRedis.On("Set", mock.Anything, "session:"+refreshToken, mock.Anything, sessionExpiry).Return(nil)
	mockRedis.On("Get", mock.Anything, "session:"+refreshToken).Return(nil)

	// Store session
	err := service.storeSession(context.Background(), refreshToken, sessionInfo)
	assert.NoError(t, err)

	// Get session
	retrievedSession, err := service.getSession(context.Background(), refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, sessionInfo.UserID, retrievedSession.UserID)
	assert.Equal(t, sessionInfo.Email, retrievedSession.Email)
	assert.Equal(t, sessionInfo.Name, retrievedSession.Name)
}
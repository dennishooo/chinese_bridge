package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"chinese-bridge-game/internal/auth/dto"
	"chinese-bridge-game/internal/auth/repository"
	"chinese-bridge-game/internal/common/config"
	"chinese-bridge-game/internal/common/database"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

const (
	// Redis key prefixes
	sessionPrefix = "session:"
	userPrefix    = "user:"
	
	// Token expiration times
	accessTokenExpiry  = 1 * time.Hour
	refreshTokenExpiry = 24 * time.Hour * 7 // 7 days
	sessionExpiry      = 24 * time.Hour * 7 // 7 days
)

type AuthService interface {
	GoogleOAuthLogin(ctx context.Context, code string) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*dto.JWTClaims, error)
	Logout(ctx context.Context, userID string) error
	GetGoogleOAuthURL(state string) string
}

type authService struct {
	repo         repository.AuthRepository
	redisClient  RedisClient
	config       *config.Config
	oauthConfig  *oauth2.Config
}

func NewAuthService(repo repository.AuthRepository, redisClient RedisClient, config *config.Config) AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     config.GoogleOAuth.ClientID,
		ClientSecret: config.GoogleOAuth.ClientSecret,
		RedirectURL:  config.GoogleOAuth.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &authService{
		repo:        repo,
		redisClient: redisClient,
		config:      config,
		oauthConfig: oauthConfig,
	}
}

func (s *authService) GetGoogleOAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *authService) GoogleOAuthLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	oauth2Service, err := oauth2v2.NewService(ctx, option.WithTokenSource(s.oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists or create new user
	user, err := s.repo.GetUserByGoogleID(ctx, userInfo.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by google id: %w", err)
	}

	if user == nil {
		// Create new user
		user = &database.User{
			ID:       uuid.New().String(),
			GoogleID: userInfo.Id,
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Avatar:   userInfo.Picture,
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		// Create initial user stats
		stats := &database.UserStats{
			UserID: user.ID,
		}
		user.Stats = stats
	} else {
		// Update existing user info
		user.Name = userInfo.Name
		user.Avatar = userInfo.Picture
		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Generate JWT tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store session in Redis
	sessionInfo := &dto.SessionInfo{
		UserID:       user.ID,
		Email:        user.Email,
		Name:         user.Name,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(sessionExpiry),
	}

	if err := s.storeSession(ctx, refreshToken, sessionInfo); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Store session in database
	dbSession := &database.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: sessionInfo.ExpiresAt,
	}

	if err := s.repo.CreateSession(ctx, dbSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenExpiry.Seconds()),
		User: dto.UserInfo{
			ID:       user.ID,
			GoogleID: user.GoogleID,
			Email:    user.Email,
			Name:     user.Name,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	// Get session from Redis
	sessionInfo, err := s.getSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if sessionInfo == nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if session is expired
	if time.Now().After(sessionInfo.ExpiresAt) {
		// Clean up expired session
		s.deleteSession(ctx, refreshToken)
		s.repo.DeleteSession(ctx, refreshToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user from database
	user, err := s.repo.GetUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenExpiry.Seconds()),
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*dto.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email claim")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name claim")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid iat claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp claim")
	}

	return &dto.JWTClaims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// Delete all user sessions from database
	if err := s.repo.DeleteUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	// Delete user sessions from Redis (this is a simplified approach)
	// In a production system, you might want to maintain a mapping of user to sessions
	pattern := sessionPrefix + "*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get session keys: %w", err)
	}

	for _, key := range keys {
		sessionData, err := s.redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionInfo dto.SessionInfo
		if err := json.Unmarshal([]byte(sessionData), &sessionInfo); err != nil {
			continue
		}

		if sessionInfo.UserID == userID {
			s.redisClient.Del(ctx, key)
		}
	}

	return nil
}

func (s *authService) generateAccessToken(user *database.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"iat":     now.Unix(),
		"exp":     now.Add(accessTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *authService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *authService) storeSession(ctx context.Context, refreshToken string, sessionInfo *dto.SessionInfo) error {
	sessionData, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}

	key := sessionPrefix + refreshToken
	return s.redisClient.Set(ctx, key, sessionData, sessionExpiry).Err()
}

func (s *authService) getSession(ctx context.Context, refreshToken string) (*dto.SessionInfo, error) {
	key := sessionPrefix + refreshToken
	sessionData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var sessionInfo dto.SessionInfo
	if err := json.Unmarshal([]byte(sessionData), &sessionInfo); err != nil {
		return nil, err
	}

	return &sessionInfo, nil
}

func (s *authService) deleteSession(ctx context.Context, refreshToken string) error {
	key := sessionPrefix + refreshToken
	return s.redisClient.Del(ctx, key).Err()
}
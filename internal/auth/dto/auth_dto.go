package dto

import "time"

// GoogleOAuthRequest represents the request for Google OAuth login
type GoogleOAuthRequest struct {
	Code  string `json:"code" binding:"required" example:"4/0AX4XfWjYZ..."`
	State string `json:"state,omitempty" example:"random_state_string"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresIn    int       `json:"expires_in" example:"3600"`
	User         UserInfo  `json:"user"`
}

// TokenResponse represents the response for token refresh
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

// RefreshTokenRequest represents the request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	GoogleID string `json:"google_id" example:"1234567890"`
	Email    string `json:"email" example:"user@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Avatar   string `json:"avatar" example:"https://lh3.googleusercontent.com/..."`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"Invalid request parameters"`
	Details string `json:"details,omitempty" example:"Field 'code' is required"`
	TraceID string `json:"trace_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// GoogleUserInfo represents user info from Google OAuth API
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IssuedAt int64  `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

// SessionInfo represents session information stored in Redis
type SessionInfo struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
package middleware

import (
	"net/http"
	"strings"
	"time"

	"chinese-bridge-game/internal/auth/dto"
	"chinese-bridge-game/internal/auth/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// JWTAuth middleware for JWT token validation
func JWTAuth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Authorization header is required",
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Invalid authorization header format",
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Token is required",
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Invalid or expired token",
				Details: err.Error(),
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)

		c.Next()
	}
}

// RateLimiter middleware for rate limiting
func RateLimiter(requestsPerSecond float64, burstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burstSize)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Too many requests",
				Details: "Rate limit exceeded, please try again later",
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimiter middleware for per-IP rate limiting
func IPRateLimiter(requestsPerSecond float64, burstSize int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burstSize)
			limiters[ip] = limiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Too many requests from this IP",
				Details: "Rate limit exceeded, please try again later",
				TraceID: c.GetString("trace_id"),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// TraceID middleware adds a trace ID to each request
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// Generate a simple trace ID (in production, use a proper UUID library)
			traceID = generateTraceID()
		}
		
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		
		c.Next()
	}
}

// generateTraceID generates a simple trace ID
func generateTraceID() string {
	return time.Now().Format("20060102150405") + "-" + strings.Replace(time.Now().Format("000000"), "0", "", -1)
}
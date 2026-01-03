package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/config"
	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/internal/repository"
)

const (
	// UserContextKey is the key for storing user info in Fiber context
	UserContextKey = "user"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(userRepo repository.UserRepository, cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// RequireAuth requires a valid JWT token
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Authorization header required"))
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid authorization format"))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Token required"))
		}

		// Parse and validate token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			log.Debug().Err(err).Msg("Token validation failed")
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid or expired token"))
		}

		// Get user from database
		user, err := m.userRepo.FindByID(c.Context(), claims.UserID)
		if err != nil {
			log.Debug().Err(err).Str("user_id", claims.UserID).Msg("User not found")
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("User not found"))
		}

		// Store user in context (without password)
		user.Password = nil
		c.Locals(UserContextKey, user)

		return c.Next()
	}
}

// RequireAdmin requires the user to have admin role
func (m *AuthMiddleware) RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals(UserContextKey).(*domain.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Authentication required"))
		}

		if !user.IsAdmin() {
			log.Debug().
				Str("user_id", user.ID).
				Str("role", user.Role).
				Msg("User is not admin")
			return c.Status(fiber.StatusForbidden).JSON(domain.ErrorResponse("Admin access required"))
		}

		return c.Next()
	}
}

// OptionalAuth tries to authenticate but doesn't require it
func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next()
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Next()
		}

		claims, err := m.validateToken(tokenString)
		if err != nil {
			return c.Next()
		}

		user, err := m.userRepo.FindByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Next()
		}

		user.Password = nil
		c.Locals(UserContextKey, user)

		return c.Next()
	}
}

// validateToken validates a JWT token and returns its claims
func (m *AuthMiddleware) validateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// GetUserFromContext gets the user from the Fiber context
func GetUserFromContext(c *fiber.Ctx) (*domain.User, bool) {
	user, ok := c.Locals(UserContextKey).(*domain.User)
	return user, ok
}

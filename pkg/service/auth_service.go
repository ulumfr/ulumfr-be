package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/ulumfr/ulumfr-be/pkg/config"
	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	r2Client    *storage.R2Client
	cfg         *config.Config
	validate    *validator.Validate
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, r2Client *storage.R2Client, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		r2Client:    r2Client,
		cfg:         cfg,
		validate:    validator.New(),
	}
}

// Register handles user registration
func (s *AuthService) Register(c *fiber.Ctx) error {
	var input domain.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(c.Context(), input.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(domain.ErrorResponse("User with this email already exists"))
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to register user"))
	}

	// Create user
	user, err := s.userRepo.Create(c.Context(), input.Name, input.Email, string(hashedPassword))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to register user"))
	}

	// Generate tokens and create session
	tokens, err := s.generateTokensWithSession(c, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to register user"))
	}

	log.Info().Str("email", input.Email).Msg("User registered successfully")

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(tokens, "Registration successful"))
}

// Login handles user login
func (s *AuthService) Login(c *fiber.Ctx) error {
	var input domain.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(c.Context(), input.Email)
	if err != nil {
		log.Debug().Str("email", input.Email).Msg("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid email or password"))
	}

	// Check if user has password (local auth)
	if user.Password == nil || *user.Password == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid email or password"))
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		log.Debug().Str("email", input.Email).Msg("Invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid email or password"))
	}

	// Generate tokens and create session
	tokens, err := s.generateTokensWithSession(c, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to login"))
	}

	log.Info().Str("email", input.Email).Msg("User logged in successfully")

	return c.JSON(domain.SuccessResponse(tokens, "Login successful"))
}

// Logout handles user logout
func (s *AuthService) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Not authenticated"))
	}

	var input domain.LogoutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("refresh_token is required"))
	}

	// Verify the refresh token belongs to this user
	session, err := s.sessionRepo.FindByToken(c.Context(), input.RefreshToken)
	if err != nil || session.UserID != user.ID {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid refresh token"))
	}

	// Delete the session
	if err := s.sessionRepo.Delete(c.Context(), input.RefreshToken); err != nil {
		log.Error().Err(err).Msg("Failed to delete session")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to logout"))
	}

	log.Info().Str("user_id", user.ID).Msg("User logged out successfully")

	return c.JSON(domain.SuccessResponse(nil, "Logout successful"))
}

// LogoutAll handles logging out from all devices
func (s *AuthService) LogoutAll(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Not authenticated"))
	}

	// Delete all sessions for this user
	if err := s.sessionRepo.DeleteByUserID(c.Context(), user.ID); err != nil {
		log.Error().Err(err).Msg("Failed to delete all sessions")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to logout from all devices"))
	}

	log.Info().Str("user_id", user.ID).Msg("User logged out from all devices")

	return c.JSON(domain.SuccessResponse(nil, "Logged out from all devices"))
}

// RefreshToken handles token refresh
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	var input domain.RefreshTokenInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	// Check if refresh token exists in session table
	session, err := s.sessionRepo.FindByToken(c.Context(), input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Invalid or expired refresh token"))
	}

	// Check if session is expired
	if session.IsExpired() {
		// Delete expired session
		_ = s.sessionRepo.Delete(c.Context(), input.RefreshToken)
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Refresh token expired"))
	}

	// Get user from database
	user, err := s.userRepo.FindByID(c.Context(), session.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("User not found"))
	}

	// Delete old session
	_ = s.sessionRepo.Delete(c.Context(), input.RefreshToken)

	// Generate new tokens and create new session
	tokens, err := s.generateTokensWithSession(c, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to refresh token"))
	}

	return c.JSON(domain.SuccessResponse(tokens, "Token refreshed successfully"))
}

// Me returns the current authenticated user
func (s *AuthService) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Not authenticated"))
	}

	// Remove password from response
	user.Password = nil

	return c.JSON(domain.SuccessResponse(user, ""))
}

// UpdateProfile updates the current user's profile
func (s *AuthService) UpdateProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse("Not authenticated"))
	}

	var input domain.UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	// If changing password, verify current password
	if input.NewPassword != nil {
		if input.CurrentPassword == nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Current password is required to change password"))
		}

		// Get user with password from database
		dbUser, err := s.userRepo.FindByID(c.Context(), user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch user"))
		}

		if dbUser.Password == nil || *dbUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Cannot change password for OAuth users"))
		}

		// Verify current password
		if err := bcrypt.CompareHashAndPassword([]byte(*dbUser.Password), []byte(*input.CurrentPassword)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Current password is incorrect"))
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("Failed to hash new password")
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update password"))
		}
		hashedPasswordStr := string(hashedPassword)
		input.NewPassword = &hashedPasswordStr
	}

	// If changing email, check if new email already exists
	if input.Email != nil && *input.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(c.Context(), *input.Email)
		if existingUser != nil {
			return c.Status(fiber.StatusConflict).JSON(domain.ErrorResponse("Email already in use"))
		}
	}

	// If changing image, delete old image from R2
	if input.Image != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		dbUser, err := s.userRepo.FindByID(c.Context(), user.ID)
		if err == nil && dbUser.Image != nil && *dbUser.Image != *input.Image {
			key := storage.ExtractKeyFromURL(*dbUser.Image)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old profile image from R2")
				}
			}
		}
	}

	// Update user profile
	updatedUser, err := s.userRepo.Update(c.Context(), user.ID, input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user profile")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update profile"))
	}

	// Remove password from response
	updatedUser.Password = nil

	log.Info().Str("user_id", user.ID).Msg("User profile updated successfully")

	return c.JSON(domain.SuccessResponse(updatedUser, "Profile updated successfully"))
}

// ListUsers returns all users (admin endpoint)
func (s *AuthService) ListUsers(c *fiber.Ctx) error {
	users, err := s.userRepo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch users")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch users"))
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = nil
	}

	return c.JSON(domain.SuccessResponse(users, ""))
}

// generateTokensWithSession generates tokens and stores refresh token in sessions table
func (s *AuthService) generateTokensWithSession(c *fiber.Ctx, user *domain.User) (*domain.TokenResponse, error) {
	// Create access token (JWT)
	accessClaims := domain.NewJWTClaims(user.ID, user.Email, user.Role, s.cfg.JWTAccessExpiry)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token (UUID stored in sessions table)
	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(s.cfg.JWTRefreshExpiry)

	// Store refresh token in sessions table
	_, err = s.sessionRepo.Create(c.Context(), user.ID, refreshToken, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

// validateToken validates a JWT token and returns its claims
func (s *AuthService) validateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

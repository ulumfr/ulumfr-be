package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginInput is the input for user login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterInput is the input for user registration
type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// TokenResponse is the response containing JWT tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// RefreshTokenInput is the input for refreshing tokens
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// JWTClaims represents the claims stored in JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTClaims creates new JWT claims
func NewJWTClaims(userID, email, role string, expiry time.Duration) *JWTClaims {
	now := time.Now()
	return &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "ulumfr-be",
			Subject:   userID,
		},
	}
}

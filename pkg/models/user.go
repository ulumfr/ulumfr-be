package models

import "time"

// User represents a user in the system
type User struct {
	ID            string     `json:"id"`
	Name          *string    `json:"name,omitempty"`
	Email         string     `json:"email"`
	Password      *string    `json:"-"` // Never expose password in JSON
	EmailVerified *time.Time `json:"email_verified,omitempty"`
	Image         *string    `json:"image,omitempty"`
	Role          string     `json:"role"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Session represents an active user session
type Session struct {
	ID           string    `json:"id"`
	SessionToken string    `json:"session_token"`
	UserID       string    `json:"user_id"`
	Expires      time.Time `json:"expires"`
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "ADMIN"
}

// UpdateProfileInput is the input for updating user profile
type UpdateProfileInput struct {
	Name            *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email           *string `json:"email,omitempty" validate:"omitempty,email"`
	CurrentPassword *string `json:"current_password,omitempty"`
	NewPassword     *string `json:"new_password,omitempty" validate:"omitempty,min=6,max=100"`
	Image           *string `json:"image,omitempty"`
}

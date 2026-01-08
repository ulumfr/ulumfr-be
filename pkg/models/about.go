package models

import "time"

// About represents the about me / profile information
type About struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Nickname  *string   `json:"nickname,omitempty"`
	Role      string    `json:"role"`
	Bio       *string   `json:"bio,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CoverURL  *string   `json:"cover_url,omitempty"`
	Location  *string   `json:"location,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAboutInput is the input for creating an about entry
type CreateAboutInput struct {
	FullName  string  `json:"full_name" validate:"required,min=1,max=255"`
	Nickname  *string `json:"nickname" validate:"omitempty,max=100"`
	Role      string  `json:"role" validate:"required,min=1,max=255"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
	CoverURL  *string `json:"cover_url" validate:"omitempty,url"`
	Location  *string `json:"location" validate:"omitempty,max=255"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Phone     *string `json:"phone" validate:"omitempty,max=20"`
	IsActive  bool    `json:"is_active"`
}

// UpdateAboutInput is the input for updating an about entry
type UpdateAboutInput struct {
	FullName  *string `json:"full_name" validate:"omitempty,min=1,max=255"`
	Nickname  *string `json:"nickname" validate:"omitempty,max=100"`
	Role      *string `json:"role" validate:"omitempty,min=1,max=255"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
	CoverURL  *string `json:"cover_url" validate:"omitempty,url"`
	Location  *string `json:"location" validate:"omitempty,max=255"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Phone     *string `json:"phone" validate:"omitempty,max=20"`
	IsActive  *bool   `json:"is_active"`
}

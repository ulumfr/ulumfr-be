package models

import "time"

// Contact represents a contact form submission
type Contact struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Subject   *string   `json:"subject,omitempty"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// Resume represents a resume/CV file
type Resume struct {
	ID        string    `json:"id"`
	FileURL   string    `json:"file_url"`
	FileName  string    `json:"file_name"`
	FileSize  *int      `json:"file_size,omitempty"`
	Version   *string   `json:"version,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateContactInput is the input for creating a contact submission
type CreateContactInput struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Email   string  `json:"email" validate:"required,email"`
	Subject *string `json:"subject" validate:"omitempty,max=255"`
	Message string  `json:"message" validate:"required,max=5000"`
}

// CreateResumeInput is the input for creating a resume entry
type CreateResumeInput struct {
	FileURL  string  `json:"file_url" validate:"required,url"`
	FileName string  `json:"file_name" validate:"required,min=1,max=255"`
	FileSize *int    `json:"file_size"`
	Version  *string `json:"version" validate:"omitempty,max=50"`
	IsActive bool    `json:"is_active"`
}

// UpdateResumeInput is the input for updating a resume entry
type UpdateResumeInput struct {
	FileURL  *string `json:"file_url" validate:"omitempty,url"`
	FileName *string `json:"file_name" validate:"omitempty,min=1,max=255"`
	FileSize *int    `json:"file_size"`
	Version  *string `json:"version" validate:"omitempty,max=50"`
	IsActive *bool   `json:"is_active"`
}

// ContactListParams contains parameters for listing contacts
type ContactListParams struct {
	Page   int   `query:"page"`
	Limit  int   `query:"limit"`
	IsRead *bool `query:"is_read"`
}

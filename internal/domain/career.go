package domain

import "time"

// Career represents a work experience entry
type Career struct {
	ID          string     `json:"id"`
	Company     string     `json:"company"`
	Position    string     `json:"position"`
	Location    *string    `json:"location,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsCurrent   bool       `json:"is_current"`
	LogoURL     *string    `json:"logo_url,omitempty"`
	CompanyURL  *string    `json:"company_url,omitempty"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Education represents an education entry
type Education struct {
	ID          string     `json:"id"`
	School      string     `json:"school"`
	Degree      string     `json:"degree"`
	Field       *string    `json:"field,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	GPA         *string    `json:"gpa,omitempty"`
	LogoURL     *string    `json:"logo_url,omitempty"`
	SchoolURL   *string    `json:"school_url,omitempty"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateCareerInput is the input for creating a career entry
type CreateCareerInput struct {
	Company     string     `json:"company" validate:"required,min=1,max=255"`
	Position    string     `json:"position" validate:"required,min=1,max=255"`
	Location    *string    `json:"location" validate:"omitempty,max=255"`
	Description *string    `json:"description"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   bool       `json:"is_current"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	CompanyURL  *string    `json:"company_url" validate:"omitempty,url"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateCareerInput is the input for updating a career entry
type UpdateCareerInput struct {
	Company     *string    `json:"company" validate:"omitempty,min=1,max=255"`
	Position    *string    `json:"position" validate:"omitempty,min=1,max=255"`
	Location    *string    `json:"location" validate:"omitempty,max=255"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   *bool      `json:"is_current"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	CompanyURL  *string    `json:"company_url" validate:"omitempty,url"`
	SortOrder   *int       `json:"sort_order"`
}

// CreateEducationInput is the input for creating an education entry
type CreateEducationInput struct {
	School      string     `json:"school" validate:"required,min=1,max=255"`
	Degree      string     `json:"degree" validate:"required,min=1,max=255"`
	Field       *string    `json:"field" validate:"omitempty,max=255"`
	Location    *string    `json:"location" validate:"omitempty,max=255"`
	Description *string    `json:"description"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"`
	GPA         *string    `json:"gpa" validate:"omitempty,max=20"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	SchoolURL   *string    `json:"school_url" validate:"omitempty,url"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateEducationInput is the input for updating an education entry
type UpdateEducationInput struct {
	School      *string    `json:"school" validate:"omitempty,min=1,max=255"`
	Degree      *string    `json:"degree" validate:"omitempty,min=1,max=255"`
	Field       *string    `json:"field" validate:"omitempty,max=255"`
	Location    *string    `json:"location" validate:"omitempty,max=255"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	GPA         *string    `json:"gpa" validate:"omitempty,max=20"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	SchoolURL   *string    `json:"school_url" validate:"omitempty,url"`
	SortOrder   *int       `json:"sort_order"`
}

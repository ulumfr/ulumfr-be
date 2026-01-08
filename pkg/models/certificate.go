package models

import "time"

// Certificate represents a certificate or achievement
type Certificate struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Issuer        string     `json:"issuer"`
	IssueDate     time.Time  `json:"issue_date"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	CredentialID  *string    `json:"credential_id,omitempty"`
	CredentialURL *string    `json:"credential_url,omitempty"`
	ImageURL      *string    `json:"image_url,omitempty"`
	Description   *string    `json:"description,omitempty"`
	SortOrder     int        `json:"sort_order"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// CreateCertificateInput is the input for creating a certificate
type CreateCertificateInput struct {
	Name          string     `json:"name" validate:"required,min=1,max=255"`
	Issuer        string     `json:"issuer" validate:"required,min=1,max=255"`
	IssueDate     time.Time  `json:"issue_date" validate:"required"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	CredentialID  *string    `json:"credential_id" validate:"omitempty,max=255"`
	CredentialURL *string    `json:"credential_url" validate:"omitempty,url"`
	ImageURL      *string    `json:"image_url" validate:"omitempty,url"`
	Description   *string    `json:"description"`
	SortOrder     int        `json:"sort_order"`
}

// UpdateCertificateInput is the input for updating a certificate
type UpdateCertificateInput struct {
	Name          *string    `json:"name" validate:"omitempty,min=1,max=255"`
	Issuer        *string    `json:"issuer" validate:"omitempty,min=1,max=255"`
	IssueDate     *time.Time `json:"issue_date"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	CredentialID  *string    `json:"credential_id" validate:"omitempty,max=255"`
	CredentialURL *string    `json:"credential_url" validate:"omitempty,url"`
	ImageURL      *string    `json:"image_url" validate:"omitempty,url"`
	Description   *string    `json:"description"`
	SortOrder     *int       `json:"sort_order"`
}

package domain

import "time"

// Project represents a portfolio project
type Project struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	Content      *string   `json:"content,omitempty"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	DemoURL      *string   `json:"demo_url,omitempty"`
	RepoURL      *string   `json:"repo_url,omitempty"`
	IsPublished  bool      `json:"is_published"`
	IsFeatured   bool      `json:"is_featured"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Categories []Category `json:"categories,omitempty"`
	Tags       []Tag      `json:"tags,omitempty"`
	Images     []ProjectImage `json:"images,omitempty"`
}

// ProjectImage represents an image in a project gallery
type ProjectImage struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	URL       string    `json:"url"`
	Alt       *string   `json:"alt,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// Category represents a project category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a project tag
type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProjectInput is the input for creating a project
type CreateProjectInput struct {
	Title        string   `json:"title" validate:"required,min=1,max=255"`
	Slug         string   `json:"slug" validate:"required,min=1,max=255"`
	Description  *string  `json:"description" validate:"omitempty,max=1000"`
	Content      *string  `json:"content"`
	ThumbnailURL *string  `json:"thumbnail_url" validate:"omitempty,url"`
	DemoURL      *string  `json:"demo_url" validate:"omitempty,url"`
	RepoURL      *string  `json:"repo_url" validate:"omitempty,url"`
	IsPublished  bool     `json:"is_published"`
	IsFeatured   bool     `json:"is_featured"`
	SortOrder    int      `json:"sort_order"`
	CategoryIDs  []string `json:"category_ids"`
	TagIDs       []string `json:"tag_ids"`
}

// UpdateProjectInput is the input for updating a project
type UpdateProjectInput struct {
	Title        *string  `json:"title" validate:"omitempty,min=1,max=255"`
	Slug         *string  `json:"slug" validate:"omitempty,min=1,max=255"`
	Description  *string  `json:"description" validate:"omitempty,max=1000"`
	Content      *string  `json:"content"`
	ThumbnailURL *string  `json:"thumbnail_url" validate:"omitempty,url"`
	DemoURL      *string  `json:"demo_url" validate:"omitempty,url"`
	RepoURL      *string  `json:"repo_url" validate:"omitempty,url"`
	IsPublished  *bool    `json:"is_published"`
	IsFeatured   *bool    `json:"is_featured"`
	SortOrder    *int     `json:"sort_order"`
	CategoryIDs  []string `json:"category_ids"`
	TagIDs       []string `json:"tag_ids"`
}

// CreateCategoryInput is the input for creating a category
type CreateCategoryInput struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Slug        string  `json:"slug" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateCategoryInput is the input for updating a category
type UpdateCategoryInput struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Slug        *string `json:"slug" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// CreateTagInput is the input for creating a tag
type CreateTagInput struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
	Slug string `json:"slug" validate:"required,min=1,max=50"`
}

// UpdateTagInput is the input for updating a tag
type UpdateTagInput struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=50"`
	Slug *string `json:"slug" validate:"omitempty,min=1,max=50"`
}

// ProjectListParams contains parameters for listing projects
type ProjectListParams struct {
	Page        int    `query:"page"`
	Limit       int    `query:"limit"`
	CategoryID  string `query:"category_id"`
	TagID       string `query:"tag_id"`
	IsFeatured  *bool  `query:"is_featured"`
	Search      string `query:"search"`
}

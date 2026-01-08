package models

import "time"

// Blog represents a blog post
type Blog struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     *string    `json:"excerpt,omitempty"`
	Content     *string    `json:"content,omitempty"`
	CoverImage  *string    `json:"cover_image,omitempty"`
	IsPublished bool       `json:"is_published"`
	IsFeatured  bool       `json:"is_featured"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Tags []Tag `json:"tags,omitempty"`
}

// CreateBlogInput is the input for creating a blog post
type CreateBlogInput struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Slug        string     `json:"slug" validate:"required,min=1,max=255"`
	Excerpt     *string    `json:"excerpt" validate:"omitempty,max=500"`
	Content     *string    `json:"content"`
	CoverImage  *string    `json:"cover_image" validate:"omitempty,url"`
	IsPublished bool       `json:"is_published"`
	IsFeatured  bool       `json:"is_featured"`
	PublishedAt *time.Time `json:"published_at"`
	SortOrder   int        `json:"sort_order"`
	TagIDs      []string   `json:"tag_ids"`
}

// UpdateBlogInput is the input for updating a blog post
type UpdateBlogInput struct {
	Title       *string    `json:"title" validate:"omitempty,min=1,max=255"`
	Slug        *string    `json:"slug" validate:"omitempty,min=1,max=255"`
	Excerpt     *string    `json:"excerpt" validate:"omitempty,max=500"`
	Content     *string    `json:"content"`
	CoverImage  *string    `json:"cover_image" validate:"omitempty,url"`
	IsPublished *bool      `json:"is_published"`
	IsFeatured  *bool      `json:"is_featured"`
	PublishedAt *time.Time `json:"published_at"`
	SortOrder   *int       `json:"sort_order"`
	TagIDs      []string   `json:"tag_ids"`
}

// BlogListParams contains parameters for listing blogs
type BlogListParams struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	TagID      string `query:"tag_id"`
	IsFeatured *bool  `query:"is_featured"`
	Search     string `query:"search"`
}

package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// BlogService handles blog business logic
type BlogService struct {
	repo     repository.BlogRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewBlogService creates a new blog service
func NewBlogService(repo repository.BlogRepository, r2Client *storage.R2Client) *BlogService {
	return &BlogService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// List returns published blogs (public endpoint)
func (s *BlogService) List(c *fiber.Ctx) error {
	var params models.BlogListParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid query parameters"))
	}

	blogs, total, err := s.repo.FindAll(c.Context(), params, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch blogs")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch blogs"))
	}

	pagination := models.NewPagination(params.Page, params.Limit, total)

	return c.JSON(models.PaginatedResponse[models.Blog]{
		Data:       blogs,
		Pagination: pagination,
	})
}

// GetBySlug returns a blog by slug (public endpoint)
func (s *BlogService) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Slug is required"))
	}

	blog, err := s.repo.FindBySlug(c.Context(), slug)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("Blog not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Blog not found"))
	}

	// Check if published for public access
	if !blog.IsPublished {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Blog not found"))
	}

	return c.JSON(models.SuccessResponse(blog, ""))
}

// AdminList returns all blogs including unpublished (admin endpoint)
func (s *BlogService) AdminList(c *fiber.Ctx) error {
	var params models.BlogListParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid query parameters"))
	}

	blogs, total, err := s.repo.FindAll(c.Context(), params, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch blogs")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch blogs"))
	}

	pagination := models.NewPagination(params.Page, params.Limit, total)

	return c.JSON(models.PaginatedResponse[models.Blog]{
		Data:       blogs,
		Pagination: pagination,
	})
}

// GetByID returns a blog by ID (admin endpoint)
func (s *BlogService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	blog, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Blog not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Blog not found"))
	}

	return c.JSON(models.SuccessResponse(blog, ""))
}

// Create creates a new blog
func (s *BlogService) Create(c *fiber.Ctx) error {
	var input models.CreateBlogInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	blog, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create blog")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create blog"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(blog, "Blog created successfully"))
}

// Update updates a blog
func (s *BlogService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateBlogInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// Handle R2 cleanup for cover image
	if input.CoverImage != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		existingBlog, err := s.repo.FindByID(c.Context(), id)
		if err == nil && existingBlog.CoverImage != nil && *existingBlog.CoverImage != *input.CoverImage {
			key := storage.ExtractKeyFromURL(*existingBlog.CoverImage)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old cover image from R2")
				}
			}
		}
	}

	blog, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update blog")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update blog"))
	}

	return c.JSON(models.SuccessResponse(blog, "Blog updated successfully"))
}

// Delete deletes a blog
func (s *BlogService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch blog for R2 cleanup
	blog, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Blog not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Blog not found"))
	}

	// Delete cover image from R2
	if s.r2Client != nil && s.r2Client.IsConfigured() && blog.CoverImage != nil {
		key := storage.ExtractKeyFromURL(*blog.CoverImage)
		if key != "" {
			if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
				log.Warn().Err(err).Str("key", key).Msg("Failed to delete cover image from R2")
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete blog")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete blog"))
	}

	return c.JSON(models.SuccessResponse(nil, "Blog deleted successfully"))
}

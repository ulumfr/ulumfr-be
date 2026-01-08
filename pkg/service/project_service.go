package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// ProjectService handles project business logic
type ProjectService struct {
	repo     repository.ProjectRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewProjectService creates a new project service
func NewProjectService(repo repository.ProjectRepository, r2Client *storage.R2Client) *ProjectService {
	return &ProjectService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// List returns published projects (public endpoint)
func (s *ProjectService) List(c *fiber.Ctx) error {
	params := models.ProjectListParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid query parameters"))
	}

	projects, total, err := s.repo.FindAll(c.Context(), params, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch projects")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch projects"))
	}

	pagination := models.NewPagination(params.Page, params.Limit, total)

	return c.JSON(models.PaginatedResponse[models.Project]{
		Data:       projects,
		Pagination: pagination,
	})
}

// AdminList returns all projects (admin endpoint)
func (s *ProjectService) AdminList(c *fiber.Ctx) error {
	params := models.ProjectListParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid query parameters"))
	}

	projects, total, err := s.repo.FindAll(c.Context(), params, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch projects")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch projects"))
	}

	pagination := models.NewPagination(params.Page, params.Limit, total)

	return c.JSON(models.PaginatedResponse[models.Project]{
		Data:       projects,
		Pagination: pagination,
	})
}

// GetBySlug returns a project by slug
func (s *ProjectService) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Slug is required"))
	}

	project, err := s.repo.FindBySlug(c.Context(), slug)
	if err != nil {
		log.Debug().Err(err).Str("slug", slug).Msg("Project not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Project not found"))
	}

	// Only return published projects for public endpoint
	if !project.IsPublished {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Project not found"))
	}

	return c.JSON(models.SuccessResponse(project, ""))
}

// GetByID returns a project by ID (admin endpoint)
func (s *ProjectService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	project, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Debug().Err(err).Str("id", id).Msg("Project not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Project not found"))
	}

	return c.JSON(models.SuccessResponse(project, ""))
}

// Create creates a new project
func (s *ProjectService) Create(c *fiber.Ctx) error {
	var input models.CreateProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	project, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create project")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create project"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(project, "Project created successfully"))
}

// Update updates a project
func (s *ProjectService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// If updating thumbnailUrl, delete old thumbnail from R2
	if input.ThumbnailURL != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		existingProject, err := s.repo.FindByID(c.Context(), id)
		if err == nil && existingProject.ThumbnailURL != nil && *existingProject.ThumbnailURL != *input.ThumbnailURL {
			key := storage.ExtractKeyFromURL(*existingProject.ThumbnailURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old thumbnail from R2")
				}
			}
		}
	}

	project, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update project")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update project"))
	}

	return c.JSON(models.SuccessResponse(project, "Project updated successfully"))
}

// Delete deletes a project and its images from R2
func (s *ProjectService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch project for R2 cleanup
	project, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Project not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Project not found"))
	}

	// Delete thumbnail and all images from R2
	if s.r2Client != nil && s.r2Client.IsConfigured() {
		// Delete thumbnail
		if project.ThumbnailURL != nil {
			key := storage.ExtractKeyFromURL(*project.ThumbnailURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete thumbnail from R2")
				}
			}
		}
		// Delete project images
		for _, img := range project.Images {
			key := storage.ExtractKeyFromURL(img.URL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete project image from R2")
				}
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete project")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete project"))
	}

	return c.JSON(models.SuccessResponse(nil, "Project deleted successfully"))
}

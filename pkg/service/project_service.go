package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
)

// ProjectService handles project business logic
type ProjectService struct {
	repo     repository.ProjectRepository
	validate *validator.Validate
}

// NewProjectService creates a new project service
func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		repo:     repo,
		validate: validator.New(),
	}
}

// List returns published projects (public endpoint)
func (s *ProjectService) List(c *fiber.Ctx) error {
	params := domain.ProjectListParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid query parameters"))
	}

	projects, total, err := s.repo.FindAll(c.Context(), params, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch projects")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch projects"))
	}

	pagination := domain.NewPagination(params.Page, params.Limit, total)

	return c.JSON(domain.PaginatedResponse[domain.Project]{
		Data:       projects,
		Pagination: pagination,
	})
}

// AdminList returns all projects (admin endpoint)
func (s *ProjectService) AdminList(c *fiber.Ctx) error {
	params := domain.ProjectListParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid query parameters"))
	}

	projects, total, err := s.repo.FindAll(c.Context(), params, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch projects")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch projects"))
	}

	pagination := domain.NewPagination(params.Page, params.Limit, total)

	return c.JSON(domain.PaginatedResponse[domain.Project]{
		Data:       projects,
		Pagination: pagination,
	})
}

// GetBySlug returns a project by slug
func (s *ProjectService) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Slug is required"))
	}

	project, err := s.repo.FindBySlug(c.Context(), slug)
	if err != nil {
		log.Debug().Err(err).Str("slug", slug).Msg("Project not found")
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse("Project not found"))
	}

	// Only return published projects for public endpoint
	if !project.IsPublished {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse("Project not found"))
	}

	return c.JSON(domain.SuccessResponse(project, ""))
}

// GetByID returns a project by ID (admin endpoint)
func (s *ProjectService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	project, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Debug().Err(err).Str("id", id).Msg("Project not found")
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse("Project not found"))
	}

	return c.JSON(domain.SuccessResponse(project, ""))
}

// Create creates a new project
func (s *ProjectService) Create(c *fiber.Ctx) error {
	var input domain.CreateProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	project, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create project")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to create project"))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(project, "Project created successfully"))
}

// Update updates a project
func (s *ProjectService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	var input domain.UpdateProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	project, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update project")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update project"))
	}

	return c.JSON(domain.SuccessResponse(project, "Project updated successfully"))
}

// Delete deletes a project
func (s *ProjectService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete project")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete project"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Project deleted successfully"))
}

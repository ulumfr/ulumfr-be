package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
)

// CategoryService handles category business logic
type CategoryService struct {
	repo     repository.CategoryRepository
	validate *validator.Validate
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo:     repo,
		validate: validator.New(),
	}
}

// List returns all categories (public endpoint)
func (s *CategoryService) List(c *fiber.Ctx) error {
	categories, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch categories")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch categories"))
	}

	return c.JSON(models.SuccessResponse(categories, ""))
}

// AdminList returns all categories (admin endpoint)
func (s *CategoryService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new category
func (s *CategoryService) Create(c *fiber.Ctx) error {
	var input models.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	category, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create category"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(category, "Category created successfully"))
}

// Update updates a category
func (s *CategoryService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	category, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update category")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update category"))
	}

	return c.JSON(models.SuccessResponse(category, "Category updated successfully"))
}

// Delete deletes a category
func (s *CategoryService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete category")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete category"))
	}

	return c.JSON(models.SuccessResponse(nil, "Category deleted successfully"))
}

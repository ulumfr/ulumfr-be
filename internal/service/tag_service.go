package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/internal/repository"
)

// TagService handles tag business logic
type TagService struct {
	repo     repository.TagRepository
	validate *validator.Validate
}

// NewTagService creates a new tag service
func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{
		repo:     repo,
		validate: validator.New(),
	}
}

// List returns all tags (public endpoint)
func (s *TagService) List(c *fiber.Ctx) error {
	tags, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch tags")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch tags"))
	}

	return c.JSON(domain.SuccessResponse(tags, ""))
}

// AdminList returns all tags (admin endpoint)
func (s *TagService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new tag
func (s *TagService) Create(c *fiber.Ctx) error {
	var input domain.CreateTagInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	tag, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create tag")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to create tag"))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(tag, "Tag created successfully"))
}

// Update updates a tag
func (s *TagService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	var input domain.UpdateTagInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	tag, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update tag")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update tag"))
	}

	return c.JSON(domain.SuccessResponse(tag, "Tag updated successfully"))
}

// Delete deletes a tag
func (s *TagService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete tag")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete tag"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Tag deleted successfully"))
}

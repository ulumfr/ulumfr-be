package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/internal/repository"
)

// EducationService handles education business logic
type EducationService struct {
	repo     repository.EducationRepository
	validate *validator.Validate
}

// NewEducationService creates a new education service
func NewEducationService(repo repository.EducationRepository) *EducationService {
	return &EducationService{
		repo:     repo,
		validate: validator.New(),
	}
}

// List returns all educations (public endpoint)
func (s *EducationService) List(c *fiber.Ctx) error {
	educations, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch educations")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch educations"))
	}

	return c.JSON(domain.SuccessResponse(educations, ""))
}

// AdminList returns all educations (admin endpoint)
func (s *EducationService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new education entry
func (s *EducationService) Create(c *fiber.Ctx) error {
	var input domain.CreateEducationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	education, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create education")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to create education"))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(education, "Education created successfully"))
}

// Update updates an education entry
func (s *EducationService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	var input domain.UpdateEducationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	education, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update education")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update education"))
	}

	return c.JSON(domain.SuccessResponse(education, "Education updated successfully"))
}

// Delete deletes an education entry
func (s *EducationService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete education")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete education"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Education deleted successfully"))
}

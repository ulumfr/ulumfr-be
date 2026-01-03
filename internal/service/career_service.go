package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/internal/repository"
)

// CareerService handles career business logic
type CareerService struct {
	repo     repository.CareerRepository
	validate *validator.Validate
}

// NewCareerService creates a new career service
func NewCareerService(repo repository.CareerRepository) *CareerService {
	return &CareerService{
		repo:     repo,
		validate: validator.New(),
	}
}

// List returns all careers (public endpoint)
func (s *CareerService) List(c *fiber.Ctx) error {
	careers, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch careers")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch careers"))
	}

	return c.JSON(domain.SuccessResponse(careers, ""))
}

// AdminList returns all careers (admin endpoint)
func (s *CareerService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new career entry
func (s *CareerService) Create(c *fiber.Ctx) error {
	var input domain.CreateCareerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	career, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create career")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to create career"))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(career, "Career created successfully"))
}

// Update updates a career entry
func (s *CareerService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	var input domain.UpdateCareerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	career, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update career")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update career"))
	}

	return c.JSON(domain.SuccessResponse(career, "Career updated successfully"))
}

// Delete deletes a career entry
func (s *CareerService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete career")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete career"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Career deleted successfully"))
}

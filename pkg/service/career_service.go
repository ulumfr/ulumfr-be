package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// CareerService handles career business logic
type CareerService struct {
	repo     repository.CareerRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewCareerService creates a new career service
func NewCareerService(repo repository.CareerRepository, r2Client *storage.R2Client) *CareerService {
	return &CareerService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// List returns all careers (public endpoint)
func (s *CareerService) List(c *fiber.Ctx) error {
	careers, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch careers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch careers"))
	}

	return c.JSON(models.SuccessResponse(careers, ""))
}

// AdminList returns all careers (admin endpoint)
func (s *CareerService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new career entry
func (s *CareerService) Create(c *fiber.Ctx) error {
	var input models.CreateCareerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	career, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create career")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create career"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(career, "Career created successfully"))
}

// Update updates a career entry
func (s *CareerService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateCareerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// If updating logoUrl, delete old logo from R2
	if input.LogoURL != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		existingCareer, err := s.repo.FindByID(c.Context(), id)
		if err == nil && existingCareer.LogoURL != nil && *existingCareer.LogoURL != *input.LogoURL {
			key := storage.ExtractKeyFromURL(*existingCareer.LogoURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old logo from R2")
				}
			}
		}
	}

	career, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update career")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update career"))
	}

	return c.JSON(models.SuccessResponse(career, "Career updated successfully"))
}

// Delete deletes a career entry and its logo from R2
func (s *CareerService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch career to get logoUrl for R2 cleanup
	career, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Career not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Career not found"))
	}

	// Delete logo from R2 if exists
	if s.r2Client != nil && s.r2Client.IsConfigured() && career.LogoURL != nil {
		key := storage.ExtractKeyFromURL(*career.LogoURL)
		if key != "" {
			if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
				log.Warn().Err(err).Str("key", key).Msg("Failed to delete logo from R2")
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete career")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete career"))
	}

	return c.JSON(models.SuccessResponse(nil, "Career deleted successfully"))
}

package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// EducationService handles education business logic
type EducationService struct {
	repo     repository.EducationRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewEducationService creates a new education service
func NewEducationService(repo repository.EducationRepository, r2Client *storage.R2Client) *EducationService {
	return &EducationService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// List returns all educations (public endpoint)
func (s *EducationService) List(c *fiber.Ctx) error {
	educations, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch educations")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch educations"))
	}

	return c.JSON(models.SuccessResponse(educations, ""))
}

// AdminList returns all educations (admin endpoint)
func (s *EducationService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new education entry
func (s *EducationService) Create(c *fiber.Ctx) error {
	var input models.CreateEducationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	education, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create education")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create education"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(education, "Education created successfully"))
}

// Update updates an education entry
func (s *EducationService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateEducationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// If updating logoUrl, delete old logo from R2
	if input.LogoURL != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		existingEducation, err := s.repo.FindByID(c.Context(), id)
		if err == nil && existingEducation.LogoURL != nil && *existingEducation.LogoURL != *input.LogoURL {
			key := storage.ExtractKeyFromURL(*existingEducation.LogoURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old logo from R2")
				}
			}
		}
	}

	education, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update education")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update education"))
	}

	return c.JSON(models.SuccessResponse(education, "Education updated successfully"))
}

// Delete deletes an education entry and its logo from R2
func (s *EducationService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch education to get logoUrl for R2 cleanup
	education, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Education not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Education not found"))
	}

	// Delete logo from R2 if exists
	if s.r2Client != nil && s.r2Client.IsConfigured() && education.LogoURL != nil {
		key := storage.ExtractKeyFromURL(*education.LogoURL)
		if key != "" {
			if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
				log.Warn().Err(err).Str("key", key).Msg("Failed to delete logo from R2")
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete education")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete education"))
	}

	return c.JSON(models.SuccessResponse(nil, "Education deleted successfully"))
}

package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// AboutService handles about business logic
type AboutService struct {
	repo     repository.AboutRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewAboutService creates a new about service
func NewAboutService(repo repository.AboutRepository, r2Client *storage.R2Client) *AboutService {
	return &AboutService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// GetActive returns the active about entry (public endpoint)
func (s *AboutService) GetActive(c *fiber.Ctx) error {
	about, err := s.repo.FindActive(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch active about")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("No active profile found"))
	}

	return c.JSON(models.SuccessResponse(about, ""))
}

// AdminList returns all about entries (admin endpoint)
func (s *AboutService) AdminList(c *fiber.Ctx) error {
	abouts, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch about entries")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch about entries"))
	}

	return c.JSON(models.SuccessResponse(abouts, ""))
}

// Create creates a new about entry
func (s *AboutService) Create(c *fiber.Ctx) error {
	var input models.CreateAboutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	about, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create about")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create about"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(about, "About created successfully"))
}

// Update updates an about entry
func (s *AboutService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateAboutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// Handle R2 cleanup for avatar and cover images
	if s.r2Client != nil && s.r2Client.IsConfigured() {
		existingAbout, err := s.repo.FindByID(c.Context(), id)
		if err == nil {
			// Cleanup old avatar if updating
			if input.AvatarURL != nil && existingAbout.AvatarURL != nil && *existingAbout.AvatarURL != *input.AvatarURL {
				key := storage.ExtractKeyFromURL(*existingAbout.AvatarURL)
				if key != "" {
					if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
						log.Warn().Err(err).Str("key", key).Msg("Failed to delete old avatar from R2")
					}
				}
			}
			// Cleanup old cover if updating
			if input.CoverURL != nil && existingAbout.CoverURL != nil && *existingAbout.CoverURL != *input.CoverURL {
				key := storage.ExtractKeyFromURL(*existingAbout.CoverURL)
				if key != "" {
					if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
						log.Warn().Err(err).Str("key", key).Msg("Failed to delete old cover from R2")
					}
				}
			}
		}
	}

	about, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update about")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update about"))
	}

	return c.JSON(models.SuccessResponse(about, "About updated successfully"))
}

// Delete deletes an about entry
func (s *AboutService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch about for R2 cleanup
	about, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("About not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("About not found"))
	}

	// Delete images from R2
	if s.r2Client != nil && s.r2Client.IsConfigured() {
		if about.AvatarURL != nil {
			key := storage.ExtractKeyFromURL(*about.AvatarURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete avatar from R2")
				}
			}
		}
		if about.CoverURL != nil {
			key := storage.ExtractKeyFromURL(*about.CoverURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete cover from R2")
				}
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete about")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete about"))
	}

	return c.JSON(models.SuccessResponse(nil, "About deleted successfully"))
}

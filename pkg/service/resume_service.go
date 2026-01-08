package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// ResumeService handles resume business logic
type ResumeService struct {
	repo     repository.ResumeRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewResumeService creates a new resume service
func NewResumeService(repo repository.ResumeRepository, r2Client *storage.R2Client) *ResumeService {
	return &ResumeService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// GetActive returns the active resume (public endpoint)
func (s *ResumeService) GetActive(c *fiber.Ctx) error {
	resume, err := s.repo.FindActive(c.Context())
	if err != nil {
		log.Debug().Err(err).Msg("No active resume found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("No active resume found"))
	}

	return c.JSON(models.SuccessResponse(resume, ""))
}

// List returns all resumes (admin endpoint)
func (s *ResumeService) List(c *fiber.Ctx) error {
	resumes, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch resumes")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch resumes"))
	}

	return c.JSON(models.SuccessResponse(resumes, ""))
}

// Create creates a new resume entry
func (s *ResumeService) Create(c *fiber.Ctx) error {
	var input models.CreateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	resume, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create resume")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create resume"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(resume, "Resume created successfully"))
}

// Update updates a resume entry
func (s *ResumeService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	resume, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update resume")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update resume"))
	}

	return c.JSON(models.SuccessResponse(resume, "Resume updated successfully"))
}

// Delete deletes a resume entry and its associated R2 file
func (s *ResumeService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch the resume first to get the file URL
	resume, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Resume not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Resume not found"))
	}

	// Delete the file from R2 if r2Client is configured and file URL exists
	if s.r2Client != nil && s.r2Client.IsConfigured() && resume.FileURL != "" {
		key := storage.ExtractKeyFromURL(resume.FileURL)
		if key != "" {
			if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
				log.Warn().Err(err).Str("key", key).Msg("Failed to delete file from R2, continuing with database deletion")
			} else {
				log.Info().Str("key", key).Msg("File deleted from R2")
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete resume")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete resume"))
	}

	return c.JSON(models.SuccessResponse(nil, "Resume deleted successfully"))
}

// Activate sets a resume as the active one
func (s *ResumeService) Activate(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	if err := s.repo.SetActive(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to activate resume")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to activate resume"))
	}

	return c.JSON(models.SuccessResponse(nil, "Resume activated successfully"))
}

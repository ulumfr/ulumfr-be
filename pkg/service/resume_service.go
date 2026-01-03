package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
)

// ResumeService handles resume business logic
type ResumeService struct {
	repo     repository.ResumeRepository
	validate *validator.Validate
}

// NewResumeService creates a new resume service
func NewResumeService(repo repository.ResumeRepository) *ResumeService {
	return &ResumeService{
		repo:     repo,
		validate: validator.New(),
	}
}

// GetActive returns the active resume (public endpoint)
func (s *ResumeService) GetActive(c *fiber.Ctx) error {
	resume, err := s.repo.FindActive(c.Context())
	if err != nil {
		log.Debug().Err(err).Msg("No active resume found")
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse("No active resume found"))
	}

	return c.JSON(domain.SuccessResponse(resume, ""))
}

// List returns all resumes (admin endpoint)
func (s *ResumeService) List(c *fiber.Ctx) error {
	resumes, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch resumes")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch resumes"))
	}

	return c.JSON(domain.SuccessResponse(resumes, ""))
}

// Create creates a new resume entry
func (s *ResumeService) Create(c *fiber.Ctx) error {
	var input domain.CreateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	resume, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create resume")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to create resume"))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(resume, "Resume created successfully"))
}

// Update updates a resume entry
func (s *ResumeService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	var input domain.UpdateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	resume, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update resume")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update resume"))
	}

	return c.JSON(domain.SuccessResponse(resume, "Resume updated successfully"))
}

// Delete deletes a resume entry
func (s *ResumeService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete resume")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete resume"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Resume deleted successfully"))
}

// Activate sets a resume as the active one
func (s *ResumeService) Activate(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.SetActive(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to activate resume")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to activate resume"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Resume activated successfully"))
}

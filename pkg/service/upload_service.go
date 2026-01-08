package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// UploadService handles file upload operations
type UploadService struct {
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewUploadService creates a new upload service
func NewUploadService(r2Client *storage.R2Client) *UploadService {
	return &UploadService{
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// GetPresignedURL generates a presigned URL for file upload
func (s *UploadService) GetPresignedURL(c *fiber.Ctx) error {
	if s.r2Client == nil || !s.r2Client.IsConfigured() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.ErrorResponse("File upload service is not configured"))
	}

	var req storage.PresignedURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// Validate folder
	allowedFolders := map[string]bool{
		"projects":   true,
		"resumes":    true,
		"careers":    true,
		"educations": true,
		"general":    true,
		"profiles":   true,
	}

	if req.Folder == "" {
		req.Folder = "general"
	}

	if !allowedFolders[req.Folder] {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid folder specified"))
	}

	// Generate presigned URL
	resp, err := s.r2Client.GeneratePresignedPutURL(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate presigned URL")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate upload URL"))
	}

	log.Info().
		Str("folder", req.Folder).
		Str("file_name", req.FileName).
		Str("content_type", req.ContentType).
		Msg("Presigned URL generated")

	return c.JSON(models.SuccessResponse(resp, "Upload URL generated successfully"))
}

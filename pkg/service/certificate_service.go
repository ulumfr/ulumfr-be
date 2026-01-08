package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// CertificateService handles certificate business logic
type CertificateService struct {
	repo     repository.CertificateRepository
	r2Client *storage.R2Client
	validate *validator.Validate
}

// NewCertificateService creates a new certificate service
func NewCertificateService(repo repository.CertificateRepository, r2Client *storage.R2Client) *CertificateService {
	return &CertificateService{
		repo:     repo,
		r2Client: r2Client,
		validate: validator.New(),
	}
}

// List returns all certificates (public endpoint)
func (s *CertificateService) List(c *fiber.Ctx) error {
	certificates, err := s.repo.FindAll(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch certificates")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch certificates"))
	}

	return c.JSON(models.SuccessResponse(certificates, ""))
}

// AdminList returns all certificates (admin endpoint)
func (s *CertificateService) AdminList(c *fiber.Ctx) error {
	return s.List(c)
}

// Create creates a new certificate
func (s *CertificateService) Create(c *fiber.Ctx) error {
	var input models.CreateCertificateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	certificate, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create certificate")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create certificate"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(certificate, "Certificate created successfully"))
}

// Update updates a certificate
func (s *CertificateService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	var input models.UpdateCertificateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	// Handle R2 cleanup for certificate image
	if input.ImageURL != nil && s.r2Client != nil && s.r2Client.IsConfigured() {
		existingCert, err := s.repo.FindByID(c.Context(), id)
		if err == nil && existingCert.ImageURL != nil && *existingCert.ImageURL != *input.ImageURL {
			key := storage.ExtractKeyFromURL(*existingCert.ImageURL)
			if key != "" {
				if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
					log.Warn().Err(err).Str("key", key).Msg("Failed to delete old image from R2")
				}
			}
		}
	}

	certificate, err := s.repo.Update(c.Context(), id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update certificate")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update certificate"))
	}

	return c.JSON(models.SuccessResponse(certificate, "Certificate updated successfully"))
}

// Delete deletes a certificate
func (s *CertificateService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("ID is required"))
	}

	// Fetch certificate for R2 cleanup
	certificate, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Certificate not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Certificate not found"))
	}

	// Delete image from R2
	if s.r2Client != nil && s.r2Client.IsConfigured() && certificate.ImageURL != nil {
		key := storage.ExtractKeyFromURL(*certificate.ImageURL)
		if key != "" {
			if err := s.r2Client.DeleteObject(c.Context(), key); err != nil {
				log.Warn().Err(err).Str("key", key).Msg("Failed to delete image from R2")
			}
		}
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete certificate")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete certificate"))
	}

	return c.JSON(models.SuccessResponse(nil, "Certificate deleted successfully"))
}

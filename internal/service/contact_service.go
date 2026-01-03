package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/internal/repository"
)

// ContactService handles contact business logic
type ContactService struct {
	repo     repository.ContactRepository
	validate *validator.Validate
}

// NewContactService creates a new contact service
func NewContactService(repo repository.ContactRepository) *ContactService {
	return &ContactService{
		repo:     repo,
		validate: validator.New(),
	}
}

// Create creates a new contact submission (public endpoint with rate limiting)
func (s *ContactService) Create(c *fiber.Ctx) error {
	var input domain.CreateContactInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid request body"))
	}

	if err := s.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse(err.Error()))
	}

	contact, err := s.repo.Create(c.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create contact")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to submit contact form"))
	}

	log.Info().
		Str("email", input.Email).
		Str("name", input.Name).
		Msg("New contact form submission")

	return c.Status(fiber.StatusCreated).JSON(domain.SuccessResponse(contact, "Thank you for your message! We'll get back to you soon."))
}

// List returns all contacts (admin endpoint)
func (s *ContactService) List(c *fiber.Ctx) error {
	params := domain.ContactListParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("Invalid query parameters"))
	}

	contacts, total, err := s.repo.FindAll(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch contacts")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to fetch contacts"))
	}

	pagination := domain.NewPagination(params.Page, params.Limit, total)

	return c.JSON(domain.PaginatedResponse[domain.Contact]{
		Data:       contacts,
		Pagination: pagination,
	})
}

// GetByID returns a contact by ID (admin endpoint)
func (s *ContactService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	contact, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		log.Debug().Err(err).Str("id", id).Msg("Contact not found")
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse("Contact not found"))
	}

	return c.JSON(domain.SuccessResponse(contact, ""))
}

// MarkAsRead marks a contact as read
func (s *ContactService) MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.MarkAsRead(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to mark contact as read")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to update contact"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Contact marked as read"))
}

// Delete deletes a contact
func (s *ContactService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse("ID is required"))
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete contact")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse("Failed to delete contact"))
	}

	return c.JSON(domain.SuccessResponse(nil, "Contact deleted successfully"))
}

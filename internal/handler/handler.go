package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"message": "ulumfr-be is running",
	})
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Check for Prisma not found error
	if err == db.ErrNotFound {
		code = fiber.StatusNotFound
		message = "Resource not found"
	}

	return c.Status(code).JSON(domain.ErrorResponse(message))
}

package middleware

import (
	"regexp"
	"strings"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ValidateLogin(c *fiber.Ctx) error {
	var request models.LoginRequest

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Check required fields
	if strings.TrimSpace(request.Email) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	if strings.TrimSpace(request.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(request.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	return c.Next()
}

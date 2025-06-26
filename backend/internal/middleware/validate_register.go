package middleware

import (
	"regexp"
	"strings"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// ValidateRegister middleware for user registration validation
func ValidateRegister(c *fiber.Ctx) error {
	var request models.RegisterUserRequest

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
	if strings.TrimSpace(request.Role) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role is required",
		})
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(request.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	// Password validation
	if len(request.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Role validation
	validRoles := []string{"starthub", "investor", "donator", "collaborator"}
	roleValid := false
	for _, role := range validRoles {
		if request.Role == role {
			roleValid = true
			break
		}
	}
	if !roleValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role must be either 'starthub', 'investor', 'donator' or 'collaborator'",
		})
	}

	// If all validations pass, continue to next handler
	return c.Next()
}

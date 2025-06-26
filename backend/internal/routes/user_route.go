package routes

import (
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/ecetinerdem/starthub-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.RegisterUserRequest

		if err := c.BodyParser(&request); err != nil {
			log.Printf("❌ Could not parse registration request: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid registration data",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Printf("❌ Password hashing error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not process password",
			})
		}

		user := models.User{
			Email:    request.Email,
			Password: string(hashedPassword),
			Role:     request.Role,
		}

		query := `
		INSERT INTO users (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
		`

		err = db.QueryRow(c.Context(), query, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
		if err != nil {
			log.Printf("❌ Database error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user - email might already exist",
			})
		}

		//Generate JWT here

		token, err := utils.GenerateJWT(user)
		if err != nil {
			log.Printf("❌ Could not generate JWT: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "User created but could not generate authentication token",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
			User: models.UserResponse{
				ID:        user.ID,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
			},
			Token: token,
		})
	}
}

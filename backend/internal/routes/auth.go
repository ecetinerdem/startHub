package routes

import (
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/ecetinerdem/starthub-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.LoginRequest

		if err := c.BodyParser(&request); err != nil {
			log.Printf("❌ Could not parse login request: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid login data",
			})
		}

		//Query user from db

		var user models.User

		query := `SELECT id, email, password, role, created_at FROM users WHERE email = $1`

		err := db.QueryRow(c.Context(), query, request.Email).Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
		)

		if err != nil {
			log.Printf("❌ User not found: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		//Check password

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

		if err != nil {
			log.Printf("❌ Invalid password for user %s", user.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		//Generate JWT
		token, err := utils.GenerateJWT(user)

		if err != nil {
			log.Printf("❌ Could not generate JWT: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate authentication token",
			})
		}

		//Return user data and token

		return c.Status(fiber.StatusOK).JSON(models.AuthResponse{
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

package routes

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllStarthubs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "Select * FROM starthubs")
		if err != nil {
			log.Printf("❌ DB query failed: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch starthubs",
			})
		}
		defer rows.Close()

		var starthubs []models.StartHub

		for rows.Next() {
			var s models.StartHub
			var categoryJSON, colStarthubsJSON, colCompaniesJSON, investorsJSON, donatorsJSON []byte

			err := rows.Scan(
				&s.ID,
				&s.Name,
				&categoryJSON,
				&s.Description,
				&s.Location,
				&s.TeamSize,
				&s.URL,
				&s.Email,
				&colStarthubsJSON,
				&colCompaniesJSON,
				&investorsJSON,
				&donatorsJSON,
				&s.JoinDate,
			)

			if err != nil {
				log.Printf("❌ Row scan failed: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan starthub",
				})
			}

			// Parse JSON fields into Go slices
			json.Unmarshal(categoryJSON, &s.Category)
			json.Unmarshal(colStarthubsJSON, &s.CollaboratingStarthubs)
			json.Unmarshal(colCompaniesJSON, &s.CollaboratingCompanies)
			json.Unmarshal(investorsJSON, &s.Investors)
			json.Unmarshal(donatorsJSON, &s.Donators)

			starthubs = append(starthubs, s)
		}

		if len(starthubs) == 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "No starthubs found yet. You can add one using the POST /starthubs route.",
			})
		}

		return c.JSON(starthubs)
	}
}

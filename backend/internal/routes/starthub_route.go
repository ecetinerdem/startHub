package routes

import (
	"database/sql"
	"encoding/json"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllStarthubs(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("Select id, name, category, description, location, teamsize, url, email, collaborating_starthub, collaborating_companies, investors, donators, join_date FROM starhubs")
		if err != nil {
			return c.Status(500).SendString("Failed to Fetch starthubs")
		}
		defer rows.Close()

		var starhubs []models.StartHub

		for rows.Next() {
			var s models.StartHub
			var categoryJSON, colStarthubsJSON, colCompaniesJSON, investorsJSON, donatorsJSON []byte

			err := rows.Scan(
				&s.ID, &s.Name, &categoryJSON, &s.Description, &s.Location, &s.TeamSize,
				&s.URL, &s.Email, &colStarthubsJSON, &colCompaniesJSON, &investorsJSON, &donatorsJSON, &s.JoinDate,
			)

			if err != nil {
				return c.Status(500).SendString("Error scanning starthub row")
			}

			// Parse JSON fields into Go slices
			json.Unmarshal(categoryJSON, &s.Category)
			json.Unmarshal(colStarthubsJSON, &s.CollaboratingStarthubs)
			json.Unmarshal(colCompaniesJSON, &s.CollaboratingCompanies)
			json.Unmarshal(investorsJSON, &s.Investors)
			json.Unmarshal(donatorsJSON, &s.Donators)

			starhubs = append(starhubs, s)
		}

		return c.JSON(starhubs)
	}
}

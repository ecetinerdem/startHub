package routes

import (
	"context"
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAllStarthubs - Gets all starthubs from database
func GetAllStarthubs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Write a simple SQL query that matches your database
		query := "SELECT id, name, description, location, team_size, url, email, join_date FROM starthubs"

		// Step 2: Execute the query
		rows, err := db.Query(context.Background(), query)
		if err != nil {
			log.Printf("❌ Database error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not get starthubs from database",
			})
		}
		defer rows.Close() // Always close rows when done

		// Step 3: Create a slice to hold our results
		var starthubs []models.StartHub

		// Step 4: Loop through each row and scan the data
		for rows.Next() {
			var s models.StartHub

			err := rows.Scan(
				&s.ID,
				&s.Name,
				&s.Description,
				&s.Location,
				&s.TeamSize,
				&s.URL,
				&s.Email,
				&s.JoinDate,
			)

			if err != nil {
				log.Printf("❌ Could not read row: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Could not read data from database",
				})
			}

			// Add this starthub to our list
			starthubs = append(starthubs, s)
		}

		// Step 5: Return the results as JSON
		return c.JSON(starthubs)
	}
}

// CreateStartHub - Creates a new starthub
func CreateStartHub(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Create a variable to hold the incoming data
		var s models.StartHub

		// Step 2: Parse the JSON from the request body
		if err := c.BodyParser(&s); err != nil {
			log.Printf("❌ Could not parse request: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid data in request",
			})
		}

		// Step 3: Basic validation
		if s.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name is required",
			})
		}
		if s.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email is required",
			})
		}

		// Step 4: Insert into database (only columns that exist!)
		query := `
		INSERT INTO starthubs (name, description, location, team_size, url, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, join_date
		`

		// Step 5: Execute the query and get the generated ID and join_date back
		err := db.QueryRow(
			context.Background(),
			query,
			s.Name,
			s.Description,
			s.Location,
			s.TeamSize,
			s.URL,
			s.Email,
		).Scan(&s.ID, &s.JoinDate)

		if err != nil {
			log.Printf("❌ Could not create starthub: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not save to database",
			})
		}

		// Step 6: Return the created starthub
		return c.Status(fiber.StatusCreated).JSON(s)
	}
}

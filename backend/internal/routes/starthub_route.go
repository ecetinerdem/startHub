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

// GetStartHubByID - Gets one starthub with ID
func GetStartHubByID(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Get the ID from context parameters
		id := c.Params("id")

		// Basic ID check
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID is required",
			})
		}

		// SQL query for getting one starthub based on ID
		query := "SELECT id, name, description, location, team_size, url, email, join_date FROM starthubs WHERE id = $1"

		// Initilize a starthub model to variable
		var s models.StartHub

		// Context not sure what it is but basically satisfies the context that is queried, takes the query and id
		//Queried info scanned in to s model
		err := db.QueryRow(context.Background(), query, id).Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.Location,
			&s.TeamSize,
			&s.URL,
			&s.Email,
			&s.JoinDate,
		)

		// if error id is causing error and if no result then id not there
		if err != nil {
			log.Printf("❌ Database error for ID %s: %v", id, err)

			if err.Error() == "no rows in result set" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Starthub not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not get starthub from database",
			})
		}

		// all is good then fiber context named c turns s model into json
		return c.JSON(s)
	}
}

func GetStartHubsBySearchTerm(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		searchTerm := c.Query("name")

		if searchTerm == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Search term 'name' is required.",
			})
		}

		query := "SELECT id, name, description, location, team_size, url, email, join_date FROM starthubs WHERE name ILIKE $1"
		searchPattern := "%" + searchTerm + "%"

		rows, err := db.Query(context.Background(), query, searchPattern)
		if err != nil {
			log.Printf("❌ Database error for search '%s': %v", searchTerm, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not search starthubs",
			})
		}
		defer rows.Close()

		var starthubs []models.StartHub

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
				log.Printf("❌ Row scan error: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Could not process results",
				})
			}
			starthubs = append(starthubs, s)
		}

		// Always return 200 with consistent structure
		return c.JSON(fiber.Map{
			"search_term": searchTerm,
			"found":       len(starthubs),
			"results":     starthubs, // Empty array if no results
		})
	}
}

func CreateStartHub(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Use the proper request struct
		var req models.CreateStartHubRequest

		// Step 2: Parse the JSON from the request body
		if err := c.BodyParser(&req); err != nil {
			log.Printf("❌ Could not parse request: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid data in request",
			})
		}

		// Step 3: Basic validation
		if req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name is required",
			})
		}
		if req.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email is required",
			})
		}

		// Step 4: Start a transaction for multiple table operations
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Printf("❌ Could not start transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database transaction error",
			})
		}
		defer tx.Rollback(context.Background()) // Rollback if we don't commit

		// Step 5: Insert the starthub first
		var s models.StartHub
		query := `
		INSERT INTO starthubs (name, description, location, team_size, url, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, join_date
		`

		err = tx.QueryRow(
			context.Background(),
			query,
			req.Name,
			req.Description,
			req.Location,
			req.TeamSize,
			req.URL,
			req.Email,
		).Scan(&s.ID, &s.JoinDate)

		if err != nil {
			log.Printf("❌ Could not create starthub: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not save starthub to database",
			})
		}

		// Copy the request data to our response struct
		s.Name = req.Name
		s.Description = req.Description
		s.Location = req.Location
		s.TeamSize = req.TeamSize
		s.URL = req.URL
		s.Email = req.Email

		// Step 6: Handle categories if provided
		if len(req.Categories) > 0 {
			for _, categoryName := range req.Categories {
				if categoryName == "" {
					continue // Skip empty category names
				}

				// First, ensure the category exists (insert if not exists)
				var categoryID int
				categoryQuery := `
				INSERT INTO categories (name) 
				VALUES ($1) 
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				RETURNING id
				`

				err = tx.QueryRow(context.Background(), categoryQuery, categoryName).Scan(&categoryID)
				if err != nil {
					log.Printf("❌ Could not create/get category '%s': %v", categoryName, err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Could not process categories",
					})
				}

				// Then, link the starthub to the category
				linkQuery := `
				INSERT INTO starthub_categories (starthub_id, category_id)
				VALUES ($1, $2)
				ON CONFLICT (starthub_id, category_id) DO NOTHING
				`

				_, err = tx.Exec(context.Background(), linkQuery, s.ID, categoryID)
				if err != nil {
					log.Printf("❌ Could not link starthub to category: %v", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Could not link categories",
					})
				}
			}

			// Add categories to response
			s.Categories = req.Categories
		}

		// Step 7: Commit the transaction
		err = tx.Commit(context.Background())
		if err != nil {
			log.Printf("❌ Could not commit transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not complete starthub creation",
			})
		}

		// Step 8: Return the created starthub with categories
		return c.Status(fiber.StatusCreated).JSON(s)
	}
}

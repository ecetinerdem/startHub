package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ecetinerdem/starthub-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// getImageFromPexels fetches an image URL from Pexels based on category
func getImageFromPexels(category string) string {
	// Get API key from environment
	apiKey := os.Getenv("PEXELS_API_KEY")
	if apiKey == "" {
		log.Printf("⚠️  PEXELS_API_KEY not found in environment")
		return "" // Return empty string if no API key
	}

	// Clean up the category for search (remove spaces, make lowercase)
	searchTerm := strings.ToLower(strings.TrimSpace(category))
	if searchTerm == "" {
		searchTerm = "startup" // Default fallback
	}

	// Build the API URL
	url := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&per_page=1", searchTerm)

	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("❌ Could not create Pexels request: %v", err)
		return ""
	}

	// Add the authorization header
	req.Header.Add("Authorization", apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Pexels API request failed: %v", err)
		return ""
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		log.Printf("❌ Pexels API returned status %d", resp.StatusCode)
		return ""
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Could not read Pexels response: %v", err)
		return ""
	}

	// Parse the JSON response
	var pexelsResp models.PexelsResponse
	err = json.Unmarshal(body, &pexelsResp)
	if err != nil {
		log.Printf("❌ Could not parse Pexels response: %v", err)
		return ""
	}

	// Check if we got any photos
	if len(pexelsResp.Photos) == 0 {
		log.Printf("⚠️  No photos found for category: %s", category)
		return ""
	}

	// Return the medium image URL
	imageURL := pexelsResp.Photos[0].Src.Medium
	log.Printf("✅ Got image from Pexels for '%s': %s", category, imageURL)
	return imageURL
}

// GetAllStarthubs - Gets all starthubs from database with images
func GetAllStarthubs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Updated query to include image_url
		query := "SELECT id, name, description, location, team_size, url, email, join_date, image_url FROM starthubs"

		// Execute the query
		rows, err := db.Query(context.Background(), query)
		if err != nil {
			log.Printf("❌ Database error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not get starthubs from database",
			})
		}
		defer rows.Close()

		// Create a slice to hold our results
		var starthubs []models.StartHub

		// Loop through each row and scan the data
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
				&s.ImageURL, // Added image_url field
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

		// Return the results as JSON
		return c.JSON(starthubs)
	}
}

// GetStartHubByID - Gets one starthub with ID including image
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

		query := "SELECT id, name, description, location, team_size, url, email, join_date, image_url FROM starthubs WHERE id = $1"

		// Initialize a starthub model to variable
		var s models.StartHub

		// Execute query and scan results
		err := db.QueryRow(context.Background(), query, id).Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.Location,
			&s.TeamSize,
			&s.URL,
			&s.Email,
			&s.JoinDate,
			&s.ImageURL, // Added image_url field
		)

		// Handle errors
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

		// Return the result as JSON
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

		// Updated query to include image_url
		query := "SELECT id, name, description, location, team_size, url, email, join_date, image_url FROM starthubs WHERE name ILIKE $1"
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
				&s.ImageURL, // Added image_url field
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
			"results":     starthubs,
		})
	}
}

func CreateStartHub(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID := c.Locals("user_id").(string)
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

		// Step 4: Get image from Pexels if categories are provided
		var imageURL string
		if len(req.Categories) > 0 && req.Categories[0] != "" {
			imageURL = getImageFromPexels(req.Categories[0])
		}
		// If no image found or no categories, use a default search
		if imageURL == "" {
			imageURL = getImageFromPexels("startup")
		}

		// Step 5: Start a transaction for multiple table operations
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Printf("❌ Could not start transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database transaction error",
			})
		}
		defer tx.Rollback(context.Background()) // Rollback if we don't commit

		// Step 6: Insert the starthub first (now including image_url)
		var s models.StartHub
		query := `
		INSERT INTO starthubs (name, description, location, team_size, url, email, image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
			imageURL,
			userID, // Add created_by
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
		s.ImageURL = imageURL // Include the image URL in response

		// Step 7: Handle categories if provided (same as before)
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

		// Step 8: Commit the transaction
		err = tx.Commit(context.Background())
		if err != nil {
			log.Printf("❌ Could not commit transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not complete starthub creation",
			})
		}

		// Step 9: Return the created starthub with categories and image
		return c.Status(fiber.StatusCreated).JSON(s)
	}
}

func UpdateStartHub(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get starthub ID and user ID
		starthubID := c.Params("id")
		userID := c.Locals("user_id").(string)

		// Parse request
		var req models.CreateStartHubRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Update only if user is owner and return the updated row
		query := `
			UPDATE starthubs 
			SET name=$1, description=$2, location=$3, team_size=$4, url=$5, email=$6 
			WHERE id=$7 AND created_by=$8
			RETURNING id, name, description, location, team_size, url, email, join_date, image_url
		`

		var s models.StartHub
		err := db.QueryRow(
			context.Background(),
			query,
			req.Name,
			req.Description,
			req.Location,
			req.TeamSize,
			req.URL,
			req.Email,
			starthubID,
			userID,
		).Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.Location,
			&s.TeamSize,
			&s.URL,
			&s.Email,
			&s.JoinDate,
			&s.ImageURL,
		)

		if err != nil {
			// Handle no rows found (either doesn't exist or user isn't owner)
			if err.Error() == "no rows in result set" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Starthub not found or you're not the owner",
				})
			}

			log.Printf("❌ Database error during update: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update starthub",
			})
		}

		// Return the updated starthub
		return c.JSON(s)
	}
}

func DeleteStartHub(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get starthub ID and user ID
		starthubID := c.Params("id")
		userID := c.Locals("user_id").(string)

		// Delete only if user is owner
		query := "DELETE FROM starthubs WHERE id=$1 AND created_by=$2"
		result, err := db.Exec(
			context.Background(),
			query,
			starthubID,
			userID,
		)

		if err != nil {
			log.Printf("❌ Database error during delete: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not delete starthub",
			})
		}

		if result.RowsAffected() == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Starthub not found or you're not the owner",
			})
		}

		// Return simple success message
		return c.JSON(fiber.Map{
			"message": "Starthub deleted",
		})
	}
}

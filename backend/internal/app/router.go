package app

import (
	"github.com/ecetinerdem/starthub-backend/internal/middleware"
	"github.com/ecetinerdem/starthub-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// User
	app.Post("/sign-up", middleware.ValidateRegister, routes.RegisterUser(db))

	// Starthubs
	app.Get("/starthubs", routes.GetAllStarthubs(db))
	app.Get("/starthubs/search", routes.GetStartHubsBySearchTerm(db))
	// This route LAST in its group because matches everything then
	app.Get("/starthubs/:id", routes.GetStartHubByID(db)) // Keep as last route
	app.Post("/starthubs", routes.CreateStartHub(db))
}

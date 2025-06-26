package app

import (
	"github.com/ecetinerdem/starthub-backend/internal/middleware"
	"github.com/ecetinerdem/starthub-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Auth routes (public)
	app.Post("/sign-up", middleware.ValidateRegister, routes.RegisterUser(db))
	app.Post("/sign-in", middleware.ValidateLogin, routes.LoginUser(db))

	// Protected routes - require authentication
	api := app.Group("/api", middleware.RequireAuth)

	// Starthubs (protected)
	api.Post("/starthubs", routes.CreateStartHub(db))
	//api.Put("/starthubs/:id", routes.UpdateStartHub(db))
	//api.Delete("/starthubs/:id", routes.DeleteStartHub(db))
	// Add other protected routes here as needed
	// api.Put("/starthubs/:id", routes.UpdateStartHub(db))
	// api.Delete("/starthubs/:id", routes.DeleteStartHub(db))

	// Public starthub routes (no auth required)
	app.Get("/starthubs", routes.GetAllStarthubs(db))
	app.Get("/starthubs/search", routes.GetStartHubsBySearchTerm(db))
	app.Get("/starthubs/:id", routes.GetStartHubByID(db))
}

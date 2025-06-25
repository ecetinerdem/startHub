package app

import (
	"github.com/ecetinerdem/starthub-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupRoutes(app *fiber.App, db *pgxpool.Pool) {

	app.Get("/starthubs", routes.GetAllStarthubs(db))
	app.Get("/starthubs/:id", routes.GetStartHubByID(db))
	app.Post("/starthubs", routes.CreateStartHub(db))
}

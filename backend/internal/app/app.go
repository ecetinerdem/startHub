package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(db *pgxpool.Pool) *fiber.App {

	app := fiber.New()

	// 1. CORS MUST come FIRST (before any routes)
	// This tells the browser: "Hey, it's okay for websites to call my API"
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Allow requests from ANY website (good for development)
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false, // We don't need cookies for now
	}))

	// 2. Logger is optional but helpful - shows you what requests are coming in
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	setupRoutes(app, db)

	return app
}

package app

import (
	"github.com/gofiber/fiber/v2"
)

func Init() *fiber.App {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	return app
}

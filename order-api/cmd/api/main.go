package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "GlitchBusterOrderAPI",
		AppName:       "GlitchBusterOrderAPI v1.0.0",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("text/plain", "application/json")

		return c.SendString("Hello, World 👋!")
	})
}

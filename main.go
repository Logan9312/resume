package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Serve a simple "Hello, World!" message at the root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Serve the resume.pdf file at the /resume route
	app.Get("/resume", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.pdf") // Path to your resume.pdf
	})

	log.Fatal(app.Listen(":3000"))
}

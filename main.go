package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/resume", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.pdf")
	})

	log.Fatal(app.Listen(":3000"))
}

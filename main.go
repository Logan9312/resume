package main

import (
	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Allow requests from your frontend
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Route to serve the PDF file directly
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.pdf")
	})

	// Route to serve the PDF as PNG
	app.Get("/png", func(c *fiber.Ctx) error {
		// Generate a temporary PNG from the PDF
		outputFile := "./output.png"
		cmd := exec.Command("pdftoppm", "-png", "-singlefile", "./resume.pdf", "./output")
		err := cmd.Run()
		if err != nil {
			return c.Status(500).SendString("Failed to convert PDF to PNG: " + err.Error())
		}

		// Send the PNG file as a response
		return c.SendFile(outputFile)
	})

	log.Fatal(app.Listen(":3000"))
}

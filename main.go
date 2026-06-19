package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return sendPDF(c, "./resume.pdf", "resume.pdf")
	})

	app.Get("/resume.pdf", func(c *fiber.Ctx) error {
		return sendPDF(c, "./resume.pdf", "resume.pdf")
	})

	app.Get("/resume-2-page.pdf", func(c *fiber.Ctx) error {
		return sendPDF(c, "./2-page/resume.pdf", "resume-2-page.pdf")
	})

	app.Get("/png", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.png")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}

func sendPDF(c *fiber.Ctx, path string, filename string) error {
	c.Set("Cache-Control", "no-store, max-age=0")
	c.Set("Content-Disposition", `inline; filename="`+filename+`"`)
	c.Set("Content-Type", "application/pdf")
	c.Set("Expires", "0")
	c.Set("Pragma", "no-cache")

	return c.SendFile(path)
}

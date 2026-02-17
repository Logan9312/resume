package main

import (
	"log"
	"os"
	"strings"

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
		return c.SendFile("./resume.pdf")
	})

	app.Get("/png", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.png")
	})

	app.Get("/anonymized", func(c *fiber.Ctx) error {
		resume, err := os.ReadFile("./resume.tex")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not read resume source")
		}

		replacer := strings.NewReplacer(
			"Logan Travis", "Anonymous Candidate",
			"Edmonton, Alberta", "",
			"Edmonton, AB", "",
			"Vancouver, BC", "",
		)

		c.Type("text/plain", "utf-8")
		return c.SendString(replacer.Replace(string(resume)))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}

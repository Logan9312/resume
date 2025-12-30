package main

import (
	"bytes"
	"image/png"
	"log"
	"os"

	"github.com/gen2brain/go-fitz"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var pngData []byte

func init() {
	doc, err := fitz.New("./resume.pdf")
	if err != nil {
		log.Printf("Warning: Failed to open PDF: %v", err)
		return
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		log.Printf("Warning: Failed to render PDF: %v", err)
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Printf("Warning: Failed to encode PNG: %v", err)
		return
	}
	pngData = buf.Bytes()
	log.Println("PNG pre-generated successfully")
}

func main() {
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Route to serve the PDF file directly
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./resume.pdf")
	})

	// Route to serve the PDF as PNG
	app.Get("/png", func(c *fiber.Ctx) error {
		if pngData == nil {
			return c.Status(500).SendString("PNG not available")
		}
		c.Set("Content-Type", "image/png")
		return c.Send(pngData)
	})

	// Use PORT env variable for Railway, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

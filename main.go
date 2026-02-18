package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func resolveAssetPath(name string) (string, error) {
	candidates := []string{filepath.Join(".", name)}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, name))
	}

	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		sourceDir := filepath.Dir(sourceFile)
		candidates = append(candidates, filepath.Join(sourceDir, name))
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("asset %q not found", name)
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		resumePath, err := resolveAssetPath("resume.pdf")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not read resume PDF")
		}

		return c.SendFile(resumePath)
	})

	app.Get("/png", func(c *fiber.Ctx) error {
		resumePath, err := resolveAssetPath("resume.png")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not read resume image")
		}

		return c.SendFile(resumePath)
	})

	app.Get("/anonymized", func(c *fiber.Ctx) error {
		resumeTexPath, err := resolveAssetPath("resume.tex")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not read resume source")
		}

		resume, err := os.ReadFile(resumeTexPath)
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

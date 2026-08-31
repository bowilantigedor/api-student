package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Middleware untuk memastikan request POST, PUT, PATCH membawa header Content-Type: application/json
func requireJSON(c *fiber.Ctx) error {
	method := c.Method()
	if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch {
		contentType := c.Get("Content-Type")
		if len(contentType) < 16 || contentType[:16] != "application/json" {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(WebResponse{
				Success: false,
				Message: "Unsupported Media Type: Content-Type must be application/json",
			})
		}
	}
	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "API Students - Fiber v2",
	})

	// Global Middleware
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Route Utama / Root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Student REST API v1",
			"docs":    "/api/v1/students",
		})
	})

	// Grup API v1
	api := app.Group("/api/v1")
	studentsGroup := api.Group("/students", requireJSON)

	studentsGroup.Get("", listStudents)
	studentsGroup.Get("/:id", getStudent)
	studentsGroup.Post("", createStudent)
	studentsGroup.Put("/:id", replaceStudent)
	studentsGroup.Patch("/:id", patchStudent)
	studentsGroup.Delete("/:id", deleteStudent)

	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

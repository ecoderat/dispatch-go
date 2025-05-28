package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Register routes
	app.Get("/start", func(c *fiber.Ctx) error {
		return c.SendString("Server started")
	})

	app.Get("/stop", func(c *fiber.Ctx) error {
		return c.SendString("Server stopped")
	})

	app.Get("/sent-messages", func(c *fiber.Ctx) error {
		return c.SendString("List of sent messages")
	})

	// Start the scheduler
	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Processing all unsent messages...")
		}
	}()

	app.Listen(":3000")
}

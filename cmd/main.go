package main

import (
	"log"
	"time"

	"github.com/ecoderat/dispatch-go/internal/controller"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	ctrl := controller.NewController()

	// Register routes
	app.Get("/start", ctrl.Start)

	app.Get("/stop", ctrl.Stop)

	app.Get("/messages", ctrl.GetMessages)

	// Start the scheduler
	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Processing all unsent messages...")
		}
	}()

	app.Listen(":3000")
}

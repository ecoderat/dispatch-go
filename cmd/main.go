package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ecoderat/dispatch-go/internal/controller"
	"github.com/ecoderat/dispatch-go/internal/service"
)

func main() {
	app := fiber.New()

	messageService := service.NewMessageService()
	ctrl := controller.NewMessageController(messageService)

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

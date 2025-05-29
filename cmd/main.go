package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ecoderat/dispatch-go/internal/controller"
	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/ecoderat/dispatch-go/internal/service/message"
	"github.com/ecoderat/dispatch-go/internal/service/scheduler"
)

const (
	PostgresConnectionString = ""
	ApiURL                   = ""
)

func main() {
	app := fiber.New()
	log.Println("Starting the dispatch-go server...")

	// Setup database connection
	db := setupDatabase()

	// Initialize repository, service, and controller
	msgRepo := repository.NewMessageRepository(db)
	msgDriver := driver.NewMessageDriver(ApiURL)
	msgService := message.New(msgRepo, msgDriver)
	schedService := scheduler.New(msgService)
	ctrl := controller.NewMessageController(msgService, schedService)

	// Register routes
	app.Get("/start", ctrl.Start)
	app.Get("/stop", ctrl.Stop)
	app.Get("/messages", ctrl.GetMessages)

	// Start scheduler automatically
	if err := schedService.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	app.Listen(":3000")
}

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database successfully")

	// Migrate the schema
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Database schema migrated successfully")

	// Seed the database with initial data
	if err := db.Create(&model.Message{
		Recipient: "+90555555555",
		Content:   "Welcome to DispatchGo!",
		Status:    model.StatusPending,
	}).Error; err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeded with initial data successfully")

	return db
}

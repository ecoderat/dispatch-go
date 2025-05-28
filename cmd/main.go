package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ecoderat/dispatch-go/internal/controller"
	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/ecoderat/dispatch-go/internal/service"
	"github.com/ecoderat/dispatch-go/model"
)

const (
	PostgresConnectionString = ""
)

func main() {
	app := fiber.New()
	log.Println("Starting the dispatch-go server...")

	// Setup database connection
	db := setupDatabase()

	// Initialize repository, service, and controller
	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
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
		Content:   "Welcome to Dispatch!",
		Status:    model.StatusPending,
	}).Error; err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeded with initial data successfully")

	return db
}

package main

import (
	"context"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ecoderat/dispatch-go/internal/controller"
	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/ecoderat/dispatch-go/internal/service/message"
	"github.com/ecoderat/dispatch-go/internal/service/scheduler"
)

var (
	PostgresConnectionString string
	ApiURL                   string

	ErrDBConnection   = errors.New("failed to connect to the database")
	ErrDBMigration    = errors.New("failed to migrate database schema")
	ErrDBSeed         = errors.New("failed to seed database")
	ErrSchedulerStart = errors.New("failed to start scheduler")
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.InfoLevel)

	// Load environment variables from .env file
	if err := loadEnv(logger); err != nil {
		logger.WithError(err).Fatal("Failed to load environment variables from .env file")
	}

	// Read config from environment
	PostgresConnectionString = os.Getenv("POSTGRES_CONN_STRING")
	ApiURL = os.Getenv("API_URL")
	if PostgresConnectionString == "" || ApiURL == "" {
		logger.Fatal("POSTGRES_CONN_STRING and API_URL must be set in environment variables or .env file")
	}

	app := fiber.New()
	logger.Info("Starting the dispatch-go server...")

	// Setup database connection
	db := setupDatabase(logger)

	// Initialize repository, service, and controller
	msgRepo := repository.NewMessageRepository(db, logger)
	msgDriver := driver.NewMessageDriver(ApiURL, logger)
	msgService := message.New(msgRepo, msgDriver, logger)
	schedService := scheduler.New(msgService, logger)
	ctrl := controller.NewMessageController(msgService, schedService)

	// Register routes
	app.Get("/start", ctrl.Start)
	app.Get("/stop", ctrl.Stop)
	app.Get("/messages", ctrl.GetMessages)

	// Start scheduler automatically
	if err := schedService.Start(context.Background()); err != nil {
		logger.WithError(err).Fatal(ErrSchedulerStart)
	}

	app.Listen(":3000")
}

func setupDatabase(logger *logrus.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString))
	if err != nil {
		logger.WithError(err).Fatal(ErrDBConnection)
	}
	logger.Info("Connected to the database successfully")

	// Migrate the schema
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		logger.WithError(err).Fatal(ErrDBMigration)
	}
	logger.Info("Database schema migrated successfully")

	// Seed the database with initial data
	if err := db.Create(&model.Message{
		Recipient: "+90555555555",
		Content:   "Welcome to DispatchGo!",
		Status:    model.StatusPending,
	}).Error; err != nil {
		logger.WithError(err).Fatal(ErrDBSeed)
	}
	logger.Info("Database seeded with initial data successfully")

	return db
}

func loadEnv(logger *logrus.Logger) error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return err
	}
	logger.Info("Environment variables loaded successfully from .env file")
	return nil
}

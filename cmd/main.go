package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// Environment variables
	postgresConnectionString string
	apiURL                   string

	// Command-line flags
	fillData bool

	// Custom error messages
	ErrDBConnection    = errors.New("failed to connect to the database")
	ErrDBMigration     = errors.New("failed to migrate database schema")
	ErrDBFillDummyData = errors.New("failed to fill dummy database data")
	ErrSchedulerStart  = errors.New("failed to start scheduler")
	ErrLoadEnv         = errors.New("failed to load environment variables from .env file")
	ErrMissingEnvVars  = errors.New("required environment variables are not set")
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.InfoLevel)

	flag.BoolVar(&fillData, "fill", false, "If set, fills the database with predefined data")
	flag.Parse()

	if err := loadEnv(logger); err != nil {
		logger.WithError(err).Fatal(ErrLoadEnv)
	}

	postgresConnectionString = os.Getenv("POSTGRES_CONN_STRING")
	apiURL = os.Getenv("API_URL")
	if postgresConnectionString == "" || apiURL == "" {
		logger.Fatal(ErrMissingEnvVars, ". POSTGRES_CONN_STRING and API_URL must be set.")
	}

	app := fiber.New()
	app.Use(cors.New())

	logger.Info("Starting the dispatch-go server...")

	db, err := connectDB(postgresConnectionString, logger)
	if err != nil {
		logger.WithError(err).Fatal("Database connection failed")
	}

	if err := migrateDB(db, logger); err != nil {
		logger.WithError(err).Fatal("Database migration failed")
	}

	if fillData {
		logger.Info("The --fill flag is set. Attempting to populate database with data...")
		if err := fillDatabaseData(db, logger); err != nil {
			logger.WithError(err).Error("Failed to populate database with data. Server will continue, but DB might be empty or partially filled.")
		} else {
			logger.Info("Successfully populated database with data due to --fill flag.")
		}
	} else {
		logger.Info("The --fill flag is not set. Database will not be populated with initial/dummy data automatically.")
	}

	msgRepo := repository.NewMessageRepository(db, logger)
	msgDriver := driver.NewMessageDriver(apiURL, logger)
	msgService := message.New(msgRepo, msgDriver, logger)
	schedService := scheduler.New(msgService, logger)
	ctrl := controller.NewMessageController(msgService, schedService)

	app.Get("/start", ctrl.Start)
	app.Get("/stop", ctrl.Stop)
	app.Get("/messages", ctrl.GetMessages)

	if err := schedService.Start(context.Background()); err != nil {
		logger.WithError(err).Fatal(ErrSchedulerStart)
	}

	logger.Info("Server is listening on :3000")
	if err := app.Listen(":3000"); err != nil {
		logger.WithError(err).Fatal("Failed to start Fiber server")
	}
}

func loadEnv(logger *logrus.Logger) error {
	if err := godotenv.Load(); err != nil {
		logger.Warn("Could not load .env file. Relying on environment variables if set.")
		return nil
	}
	logger.Info("Environment variables loaded successfully from .env file")
	return nil
}

func connectDB(dsn string, logger *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Error("Database connection error")
		return nil, ErrDBConnection
	}
	logger.Info("Connected to the database successfully")
	return db, nil
}

func migrateDB(db *gorm.DB, logger *logrus.Logger) error {
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		logger.WithError(err).Error("Database migration error")
		return ErrDBMigration
	}
	logger.Info("Database schema migrated successfully")
	return nil
}

func fillDatabaseData(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Populating database with data...")

	messagesToInsert := []model.Message{
		{Recipient: "+12345678901", Content: "Dummy message 1 for testing", Status: model.StatusPending, CreatedAt: time.Now().Add(-5 * time.Minute)},
		{Recipient: "+12345678902", Content: "Another dummy message for --fill", Status: model.StatusSent, CreatedAt: time.Now().Add(-10 * time.Minute)},
		{Recipient: "+12345678903", Content: "Urgent: Fill data test", Status: model.StatusFailed, CreatedAt: time.Now().Add(-2 * time.Minute)},
		{Recipient: "+12345678904", Content: "Scheduled dummy message", Status: model.StatusPending, CreatedAt: time.Now()},
	}

	for i, msg := range messagesToInsert {
		if msg.CreatedAt.IsZero() {
			msg.CreatedAt = time.Now()
		}

		if err := db.Create(&msg).Error; err != nil {
			logger.WithError(err).Errorf("Failed to insert message #%d: %s", i+1, msg.Content)
			return fmt.Errorf("%w: failed to insert message '%s': %v", ErrDBFillDummyData, msg.Content, err)
		}
		logger.Infof("Inserted message: %s", msg.Content)
	}

	logger.Info("Database data population process completed.")
	return nil
}

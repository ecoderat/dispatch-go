# Variables
GO_APP = ./cmd/main.go
BINARY_NAME = dispatch-go
DOCKER_COMPOSE_FILE = docker-compose.yml

# Default target
.PHONY: all
all: build db-start run-fill

# Build the Go application
.PHONY: build
build:
	go build -o $(BINARY_NAME) $(GO_APP)

# Run the Go application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run the Go application with --fill flag
.PHONY: run-fill
run-fill: build
	./$(BINARY_NAME) --fill

# Clean up
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# --- Multi-process targets matching README ---
.PHONY: devstack
# App + DB (no migration)
devstack: build db-start run

.PHONY: standalone
# App only (no DB)
standalone: build run

.PHONY: db-start
# Start only the database container (without app)
db-start:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: db-stop
# Stop the database container
db-stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop db

.PHONY: db-clean
# Remove all containers and volumes
db-clean:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v

.PHONY: swagger-up
# Start Swagger UI (docs on http://localhost:8080)
swagger-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d swagger

.PHONY: swagger-down
# Stop Swagger UI
swagger-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop swagger

# Help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make all         - Build the application, start the database containers, run migrations, and run the app"
	@echo "  make devstack    - Build the app, start the database, and run the app (no migration)"
	@echo "  make standalone  - Build and run the app only (no DB)"
	@echo "  make clean       - Clean up build files and database containers"
	@echo "  make db-start    - Start only the database container (without app)"
	@echo "  make db-stop     - Stop the database container"
	@echo "  make db-clean    - Clean up database containers and volumes"
	@echo "  make swagger-up  - Start Swagger UI (docs on http://localhost:8080)"
	@echo "  make swagger-down- Stop Swagger UI"

# Variables
APP_NAME := dispatch-go # Used for native build output name
DOCKER_COMPOSE_CMD := docker-compose
DOCKER_COMPOSE_BASE_FILE := docker-compose.yml
DOCKER_COMPOSE_FILL_OVERRIDE_FILE := docker-compose.fill.yml

# Default target
.DEFAULT_GOAL := help

## --------------------------------------
## Dockerized Application Lifecycle
## --------------------------------------

# Build all service images (including the Go app)
build:
	@echo "==> Building Docker images (if needed)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) build

# Start all services (app, postgres, swagger) in DETACHED mode. App runs WITHOUT --fill.
# Logs will go to Docker's logging driver. Use 'make logs' to view them.
up:
	@echo "==> Starting all services in detached mode (app without --fill)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) up --build -d
	@echo "==> Services started. Use 'make logs' or 'make logs-app' to view application logs."

# Start all services (app, postgres, swagger) in DETACHED mode. App runs WITH --fill.
up-fill:
	@echo "==> Starting all services in detached mode (app WITH --fill)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) -f $(DOCKER_COMPOSE_FILL_OVERRIDE_FILE) up --build -d
	@echo "==> Services started. Use 'make logs' or 'make logs-app' to view application logs."

# Start all services in FOREGROUND mode. App runs WITHOUT --fill.
# Logs from all services will stream directly to your terminal. Press Ctrl+C to stop.
run:
	@echo "==> Starting all services in foreground (app without --fill)..."
	@echo "==> Logs will stream here. Press Ctrl+C to stop."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) up --build

# Start all services in FOREGROUND mode. App runs WITH --fill.
# Logs from all services will stream directly to your terminal. Press Ctrl+C to stop.
run-fill:
	@echo "==> Starting all services in foreground (app WITH --fill)..."
	@echo "==> Logs will stream here. Press Ctrl+C to stop."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) -f $(DOCKER_COMPOSE_FILL_OVERRIDE_FILE) up --build

# Stop and remove all services, volumes, and networks defined in docker-compose
down:
	@echo "==> Stopping and removing all Docker Compose services, volumes, and networks..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) down -v --remove-orphans

# Restart all services (in detached mode, app without --fill)
restart: down up

# Restart with fill (in detached mode, app with --fill)
restart-fill: down up-fill

## --------------------------------------
## Viewing Docker Logs (when services are run detached with 'make up' or 'make up-fill')
## --------------------------------------

# View and follow logs for the Go application service ('app')
logs:
	@echo "==> Tailing logs for the 'app' service (Ctrl+C to stop)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) logs -f app

# View and follow logs for ALL services
logs-all:
	@echo "==> Tailing logs for all services (Ctrl+C to stop)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) logs -f

## --------------------------------------
## Standalone Docker Service Management
## --------------------------------------

db-start:
	@echo "==> Starting PostgreSQL service in detached mode..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) up -d postgres

db-stop:
	@echo "==> Stopping PostgreSQL service..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) stop postgres

db-clean:
	@echo "==> Stopping PostgreSQL service and removing its data volume..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) stop postgres
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) rm -f postgres
	@echo "==> Removing 'db_data' volume..."
	@docker volume rm $$(basename $$(pwd))_db_data || docker volume rm db_data || echo "Volume db_data (or prefixed) not found or already removed."

swagger-up:
	@echo "==> Starting Swagger UI service in detached mode (http://localhost:8081)..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) up -d swagger

swagger-down:
	@echo "==> Stopping Swagger UI service..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE_BASE_FILE) stop swagger

## --------------------------------------
## Native Go Development (Alternative Workflow - App runs on host)
## --------------------------------------

native-build:
	@echo "==> Building Go application natively..."
	go build -o $(APP_NAME) ./cmd/main.go

# Assumes DB is running (e.g., via 'make db-start' or an external DB)
native-run: native-build
	@echo "==> Running Go application natively (ensure DB is accessible)..."
	./$(APP_NAME)

native-run-fill: native-build
	@echo "==> Running Go application natively with --fill (ensure DB is accessible)..."
	./$(APP_NAME) --fill

native-clean:
	@echo "==> Cleaning native build artifacts..."
	rm -f $(APP_NAME)
	go clean

## --------------------------------------
## Utility
## --------------------------------------

clean-all: down native-clean # Stops docker services and cleans native builds
	@echo "==> Full cleanup complete."

help:
	@echo "Available Dockerized Workflow Commands:"
	@echo "  build           - Build Docker images for all services."
	@echo "  up              - Start all services DETACHED (app without --fill). Use 'make logs'."
	@echo "  up-fill         - Start all services DETACHED (app WITH --fill). Use 'make logs'."
	@echo "  run             - Start all services FOREGROUND (app without --fill). Logs stream here."
	@echo "  run-fill        - Start all services FOREGROUND (app WITH --fill). Logs stream here."
	@echo "  down            - Stop and remove all services, volumes, and networks."
	@echo "  restart         - Restart all services (detached, app without --fill)."
	@echo "  restart-fill    - Restart all services (detached, app with --fill)."
	@echo "  logs            - Tail logs for the 'app' service (when run detached)."
	@echo "  logs-all        - Tail logs for all services (when run detached)."
	@echo "  db-start        - Start only PostgreSQL (detached)."
	@echo "  db-stop         - Stop only PostgreSQL."
	@echo "  db-clean        - Stop PostgreSQL and REMOVE its data volume."
	@echo "  swagger-up      - Start Swagger UI (detached)."
	@echo "  swagger-down    - Stop Swagger UI."
	@echo ""
	@echo "Available Native Go Development Commands (App runs on host):"
	@echo "  native-build    - Build Go app natively."
	@echo "  native-run      - Run Go app natively (DB must be accessible)."
	@echo "  native-run-fill - Run Go app natively with --fill (DB must be accessible)."
	@echo "  native-clean    - Clean native build artifacts."
	@echo ""
	@echo "Utility Commands:"
	@echo "  clean-all       - Perform 'down' (Docker cleanup) and 'native-clean'."
	@echo "  help            - Show this help message."

.PHONY: build up up-fill run run-fill down restart restart-fill logs logs-all db-start db-stop db-clean swagger-up swagger-down native-build native-run native-run-fill native-clean clean-all help
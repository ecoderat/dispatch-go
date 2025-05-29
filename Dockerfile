# Dockerfile

# --- Stage 1: Build ---
FROM golang:1.24.3-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go app
# GOOS=linux GOARCH=amd64: ensure it's built for Linux AMD64 (standard for Docker)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/dispatch-go ./cmd/main.go

# --- Stage 2: Run ---
# Use a small base image for the final application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the .env.example file (optional, for reference if needed inside container, though usually not directly used by the app)
COPY .env .

# Copy the built executable from the "builder" stage
COPY --from=builder /app/dispatch-go /app/dispatch-go

# This is informational; the actual port mapping happens in docker-compose.yml
EXPOSE 3000

# Command to run the executable
# The application will read its configuration from environment variables passed by docker-compose
ENTRYPOINT ["/app/dispatch-go"]

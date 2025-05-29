# DispatchGo

DispatchGo is a Go-based automated SMS dispatch service. It periodically sends unsent messages from a PostgreSQL database via an external provider, ensuring each is sent once. The system includes a control API, runs entirely in Docker, and is managed by a Makefile.

## Features

*   **Automated SMS Dispatch:**
    *   Periodically (e.g., every 2 minutes) retrieves unsent messages from the database.
    *   Sends messages via a configurable external SMS provider API.
*   **REST API Endpoints:**
    *   `GET /start`: Activates/re-activates the automatic message sending scheduler.
    *   `GET /stop`: Deactivates the automatic message sending scheduler.
    *   `GET /messages`: Retrieves a list of unsent messages from the database (currently lists all, future support for filtering/pagination).

## Prerequisites

*   **Docker & Docker Compose:** Essential for building and running all components of the application.
*   **Make:** For using the provided `Makefile` to simplify Docker Compose commands.
*   **Git:** For cloning the repository.
*   **(Optional for Native Go Development):** Go 1.21+ (or as specified in `go.mod`) if you choose to use the `native-*` make targets.

## Getting Started & Environment Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/ecoderat/dispatch-go.git
    cd dispatch-go
    ```

2.  **Environment Configuration (`.env` file):**
    *   Create a `.env` file in the project root by copying the sample (e.g., `.env.example` or `.env.sample`):
        ```sh
        cp .env.example .env
        ```
    *   Edit the `.env` file. **Crucially, for the Go application running in Docker to connect to the PostgreSQL container, use the Docker service name as the host:**
        ```env
        # .env
        POSTGRES_CONN_STRING= # For testing/demo, this may not need changing from the sample.
        API_URL= # For testing, see the [Webhook.site Setup](#simulating-an-sms-provider-with-webhooksite-for-developmenttesting) section below.
        ```
    *   Ensure credentials (`user`, `password`, `dbname`) in `POSTGRES_CONN_STRING` match the `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` environment variables for the `postgres` service in your `docker-compose.yml`.

3.  **(Optional, if modifying Go code) Install Go Dependencies for IDE/Native Build:**
    While the Docker build handles Go dependencies, for local Go development/IDE support:
    ```sh
    go mod tidy
    ```

## Usage with Makefile (Dockerized Workflow)

The `Makefile` primarily uses Docker Compose to build, run, and manage all application components (Go app, PostgreSQL, Swagger UI).

### Running the Application:

*   **Start all services (App, DB, Swagger) in FOREGROUND (logs stream to terminal):**
    *   Without data seeding:
        ```sh
        make run
        ```
    *   With data seeding (`--fill` flag passed to Go app):
        ```sh
        make run-fill
        ```
    Press `Ctrl+C` in the terminal to stop all services.

*   **Start all services in DETACHED mode (runs in background):**
    *   Without data seeding:
        ```sh
        make up
        ```
    *   With data seeding (`--fill` flag passed to Go app):
        ```sh
        make up-fill
        ```
    When running detached, use `make logs` or `make logs-all` to view logs.

### Stopping and Cleaning:

*   **Stop and remove all services, networks, and volumes:**
    ```sh
    make down
    ```
*   **Clean PostgreSQL data volume (DELETES ALL DB DATA):**
    ```sh
    make db-clean
    ```

### Viewing Logs (when running detached):

*   **View and follow logs for the Go application service:**
    ```sh
    make logs
    ```
*   **View and follow logs for ALL services:**
    ```sh
    make logs-all
    ```

### Building Images:
*   Build/rebuild Docker images for all services (including Go app):
    ```sh
    make build
    ```

## API Documentation (Swagger UI)

*   The OpenAPI specification (e.g., `swagger.yaml` or `docs/swagger.yml`) is served by the `swagger` service defined in `docker-compose.yml`.
*   **Start Swagger UI (if not already running via `make up` or `make run`):**
    ```sh
    make swagger-up
    ```
    Access it in your browser, typically at `http://localhost:8081` (the host port might be different based on your `docker-compose.yml`).
*   **Stop Swagger UI service:**
    ```sh
    make swagger-down
    ```

## Native Go Development (Optional Workflow)

For developers actively working on the Go code who prefer direct host execution for faster iteration or debugging:

*   **`make native-build`**: Builds the Go binary on your host.
*   **`make native-run`**: Runs the natively built Go app. *Requires the database to be accessible (e.g., started via `make db-start`)*.
*   **`make native-run-fill`**: Runs the native Go app with the `--fill` flag.
*   **`make native-clean`**: Cleans native build artifacts.

## Other Key Makefile Commands

*   **`make help`**: Displays a detailed list of all available `Makefile` commands and their descriptions for both Dockerized and Native Go workflows.

## Simulating an SMS Provider with Webhook.site (for Development/Testing)

To test the SMS sending functionality without integrating with a real SMS provider, you can use [webhook.site](https://webhook.site/). This service provides a temporary, unique URL that can receive HTTP requests, allowing you to inspect what your application sends.

**Steps to set up Webhook.site as your mock SMS API:**

1.  **Visit Webhook.site:**
    Open [https://webhook.site/](https://webhook.site/) in your web browser.
    A unique URL will be automatically generated for you (e.g., `https://webhook.site/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`).

2.  **Copy Your Unique URL:**
    This is the URL your application will send SMS API requests to. Copy it immediately.
    It will look like: `https://webhook.site/YOUR_UNIQUE_ID`

3.  **Configure the Response:**
    To make webhook.site behave more like a real API that acknowledges requests, you can customize its response:
    *   On the webhook.site page for your unique URL, find the **"Edit"** button or section (usually in the top-right area).
    *   Set the following default response parameters:
        *   **Status Code:** `202` (Accepted - common for asynchronous operations like sending an SMS)
        *   **Content type:** `application/json`
        *   **Content (Body):**
            ```json
            {
              "message": "Accepted for delivery",
              "messageId": "$request.uuid$"
            }
            ```
            *   `$request.uuid$` is a special webhook.site variable that will insert a unique ID for each request received, simulating a message ID from a provider.
    *   Save these changes on webhook.site.

4.  **Update Your `.env` File:**
    *   Open your project's `.env` file.
    *   Set the `API_URL` variable to the unique URL you copied from webhook.site:
        ```env
        # .env
        # ...
        API_URL="https://webhook.site/YOUR_UNIQUE_ID_YOU_COPIED"
        ```

## Troubleshooting

*   **PostgreSQL Connection/Authentication:**
    *   Ensure `POSTGRES_CONN_STRING` in your `.env` uses `host=postgres` (the service name) when the app runs in Docker.
    *   Verify credentials in `.env` match `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` in `docker-compose.yml` for the `postgres` service.
    *   If issues persist after configuration checks, try `make db-clean` followed by `make up` (or `make run`).
*   **Port Conflicts:** If `docker-compose up` fails, check if ports (e.g., 3000, 5432, 8081) are already in use on your host. Adjust host port mappings in `docker-compose.yml` if needed.
*   **Log Inspection:** Always use `make logs` or `docker-compose logs <service_name>` to inspect detailed error messages from containers.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

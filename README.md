# dispatch-go

## Usage Scenarios

### 1. App + DB + Migration
Run the application with the database and automatic migration:
```sh
make all
```
This will:
- Start the database container
- Run database migrations
- Build and run the Go application

### 2. App + DB (No Migration)
Run the application with the database (no migration):
```sh
make devstack
```
This will:
- Build the Go application
- Start the database container
- Run the application (without running migrations)

### 3. App Only (No DB)
Run the application only (database must already be running):
```sh
make standalone
```
This will:
- Build the Go application
- Run the application (assumes the database is already running and migrated)

### Other Commands
- **Stop Docker Containers:**
  ```sh
  make db-stop
  ```
- **Start Only Database (without app):**
  ```sh
  make db-start
  ```
- **Remove All Containers and Volumes:**
  ```sh
  make db-clean
  ```
- **Clean Up Build Artifacts:**
  ```sh
  make clean
  ```
- **Show Help:**
  ```sh
  make help
  ```

> **Note:**
> Ensure you have a valid `.env` file with the required configuration (see `.env.example` or project documentation).

This `Makefile` streamlines your build and deployment process, making it easier to manage the database and application lifecycle.
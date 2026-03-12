## Dentist API (Go + Gin)

Backend service for managing a dental clinic, built in Go using the Gin framework and following an Onion Architecture.  
It models tenants, users, dentists, patients, appointments, payments, and services, and exposes a simple users API as a starting point.

### Tech stack
- **Language**: Go
- **Web framework**: `github.com/gin-gonic/gin`
- **Database**: PostgreSQL (via GORM)
- **Architecture**: Onion (domain, application/usecase, infrastructure, API)
- **Docs**: Swagger UI (via `swaggo/gin-swagger`)

### Project structure (high level)
- **`cmd/api`**: application entrypoint (`main.go`) and composition root
- **`config`**: config loading from environment variables
- **`internal/domain`**: core entities and repository interfaces
- **`internal/application/usecase`**: business use cases orchestrating domain logic
- **`internal/infrastructure`**: DB connection and repository implementations
- **`internal/api`**: HTTP handlers, router, middleware
- **`docs`**: Swagger/OpenAPI generated files

### Running locally

#### Prerequisites
- Go installed (1.22+ recommended)
- PostgreSQL instance running and reachable

#### Environment variables
You can either create a `.env` file in the project root or export the variables directly in your shell.

Minimal variables (with defaults shown):

```bash
APP_PORT=8080           # default: 8080
APP_ENV=development    # default: development

DB_HOST=localhost      # default: localhost
DB_PORT=5432           # default: 5432
DB_USER=postgres       # default: postgres
DB_PASSWORD=postgres   # default: postgres
DB_NAME=gin_quickstart # default: gin_quickstart
DB_SSLMODE=disable     # default: disable
```

Example `.env`:

```bash
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_quickstart
DB_SSLMODE=disable
```

#### Database
Create the database in PostgreSQL (if it does not exist yet):

```bash
createdb gin_quickstart
```

Run your migrations (or let GORM auto-migrate if configured in your DB bootstrap code).

#### Start the API

```bash
go run ./cmd/api
```

The server will start on:

- **API base URL**: `http://localhost:8080/api/v1`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Health check**: `http://localhost:8080/health`

### Available endpoints (current)

#### Health
- **GET** `/health` – simple health check.

#### Users
All user routes are under `/api/v1/users`.

- **GET** `/api/v1/users` – list all users
- **GET** `/api/v1/users/:id` – get user by ID
- **POST** `/api/v1/users` – create user
- **PUT** `/api/v1/users/:id` – update user
- **DELETE** `/api/v1/users/:id` – delete user

Request/response schemas for these endpoints are documented in Swagger.

### Architecture overview

This project follows Onion Architecture:

- **Domain (`internal/domain`)**: pure domain types like `Tenants`, `Users`, `Dentists`, `Patients`, `Appointments`, `Payments`, `Services` and related enums (roles, specializations, statuses). Domain is independent of frameworks.
- **Application (`internal/application/usecase`)**: use cases such as `UserUseCase`, which depend only on domain entities and repository interfaces.
- **Infrastructure (`internal/infrastructure`)**: technical concerns like the Postgres connection and concrete repository implementations (e.g. GORM-based user repository).
- **API (`internal/api`)**: HTTP layer (Gin handlers, router, and middleware) that calls use cases and maps HTTP <-> domain.
- **Composition root (`cmd/api/main.go`)**: wires everything together (config → DB → repositories → use cases → handlers → router).

This separation keeps business rules independent from delivery and persistence concerns, making the system easier to test and evolve.

### Development tips
- **Swagger docs**: Keep handler annotations in sync with your endpoints, then regenerate Swagger docs using `swag init` (if you have it configured in your environment).
- **New features**: Add them in this order: domain (entities + repository interface) → application (use case) → infrastructure (repository impl) → API (handler + routes) → wire in `cmd/api/main.go`.
- **Configuration**: Prefer environment variables over hard-coded values; see `config/config.go` for defaults and helpers.

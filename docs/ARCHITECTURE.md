# Onion Architecture

This document describes how the **onion architecture** is applied in this project and how to work within it when adding or changing features.

## Overview

The codebase is structured in concentric layers. Dependencies point **inward**: outer layers depend on inner layers; the **domain** has no dependencies on the rest of the application.

```
                    ┌─────────────────────────────────────────┐
                    │  cmd/api/main.go (composition root)      │
                    └───────────────┬─────────────────────────┘
                                    │
         ┌──────────────────────────┼──────────────────────────┐
         │                          │                          │
         ▼                          ▼                          ▼
  config/                 internal/infrastructure/    internal/application/
  database/, repository/        usecase/
  (implements domain ports)
                                    │
                                    ▼
                          internal/domain/
                          entity/, repository/ (interfaces)
                                    ▲
                                    │
                    internal/api/ (handler → usecase only)
```

## Layers (inside-out)

### 1. Domain (`internal/domain/`)

**Role:** Core business concepts and contracts. No dependencies on other internal packages.

- **`entity/`** — Structs representing domain concepts (e.g. `User`). May include persistence tags (e.g. GORM) but no framework imports beyond stdlib and `time`.
- **`repository/`** — **Interfaces only** (ports). Define how the application reads/writes data (e.g. `UserRepository` with `FindAll`, `FindByID`, `Create`, `Update`, `Delete`). The domain owns these interfaces; infrastructure implements them.

**Rules:**

- Domain does **not** import `application`, `infrastructure`, or `api`.
- Repository interfaces use domain entities in their signatures.
- Keep entities focused on data and identity; business rules that orchestrate multiple steps belong in the application layer.

### 2. Application (`internal/application/usecase/`)

**Role:** Use cases (application logic). Orchestrate domain entities and repository ports.

- Depends **only** on `internal/domain` (entities and repository **interfaces**).
- Each use case is a struct that receives a repository interface in its constructor (e.g. `NewUserUseCase(repo repository.UserRepository)`).
- Use case methods perform one application-level operation (e.g. `GetAll`, `Create`, `Update`) and return domain entities or well-defined errors (e.g. `ErrUserNotFound`).
- Input/output can be simple DTOs (e.g. `CreateUserInput`, `UpdateUserInput`) to keep use case APIs stable.

**Rules:**

- Application does **not** import `infrastructure` or `api`.
- Application does **not** import concrete repository implementations; only `domain/repository` interfaces.
- No HTTP, no framework-specific types; only domain types and use-case-specific DTOs.

### 3. Infrastructure (`internal/infrastructure/`)

**Role:** Concrete implementations of technical concerns (database, external services).

- **`database/`** — DB connection (e.g. Postgres via GORM), migrations/auto-migrate.
- **`repository/`** — **Adapters** that implement domain repository interfaces (e.g. `GormUserRepository` implements `repository.UserRepository`). Depend on `domain/entity` and the DB client.

**Rules:**

- Infrastructure **implements** domain ports; it does not define new “domain” interfaces.
- Repository implementations live in `infrastructure/repository/` and are named to reflect the technology (e.g. `GormUserRepository`).
- Infrastructure may use `config` and third-party libraries (GORM, drivers, etc.).

### 4. API / Presentation (`internal/api/`)

**Role:** HTTP entrypoints (handlers, router, middleware).

- **`handler/`** — HTTP handlers that call use cases. Parse request (e.g. path, body), map to use-case inputs, call use case, map results/errors to HTTP status and JSON.
- **`router/`** — Registers routes and wires handlers (e.g. `/api/v1/users`).
- **`middleware/`** — Cross-cutting concerns (e.g. error recovery, logging).

**Rules:**

- API depends **only** on `internal/application/usecase` (and router/middleware). It does **not** import `domain` or `infrastructure` for business logic.
- Handlers receive use case structs via constructor (e.g. `NewUserHandler(uc *usecase.UserUseCase)`).
- Request/response DTOs (e.g. `CreateUserRequest`) live in the handler package and are mapped to/from use-case inputs and domain entities.

### 5. Composition root (`cmd/api/main.go`)

**Role:** Wire everything together.

- Load config, create DB connection, create concrete repository (e.g. `GormUserRepository`), create use case (injecting repository interface), create handler (injecting use case), create router, start server.
- Only here do concrete types from `infrastructure` and `application` meet; the rest of the app uses interfaces from the domain.

## Dependency flow summary

| Layer          | May import                          | Must not import                    |
|----------------|-------------------------------------|------------------------------------|
| Domain         | stdlib, `time`                      | application, infrastructure, api  |
| Application    | domain (entity, repository)         | infrastructure, api                |
| Infrastructure | domain, config, DB/libs             | application, api                   |
| API            | application/usecase                | domain (for business flow), infra  |
| main           | config, api, application, infra    | —                                  |

## Adding a new feature (e.g. new entity “Appointment”)

1. **Domain:** Add `internal/domain/entity/appointment.go` and `internal/domain/repository/appointment_repository.go` (interface).
2. **Application:** Add `internal/application/usecase/appointment_usecase.go` depending on `repository.AppointmentRepository` and domain entities.
3. **Infrastructure:** Add `internal/infrastructure/repository/appointment_repository.go` (e.g. `GormAppointmentRepository` implementing the interface); register migration if needed in `database` setup.
4. **API:** Add `internal/api/handler/appointment_handler.go` and register routes in `router`; request/response types in handler only.
5. **main:** In `main.go`, instantiate repo → use case → handler and pass handler (or handler group) into the router.

Use the project’s **Cursor rules** and **layer-specific skills** (domain, application, infrastructure, api) when implementing each step so the architecture stays consistent.

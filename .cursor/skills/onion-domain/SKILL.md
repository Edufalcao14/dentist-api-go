---
name: onion-domain
description: Implements the domain layer in this project's onion architecture. Use when adding or changing entities, repository interfaces (ports), or domain types in internal/domain/.
---

# Domain Layer (Onion Architecture)

Implement or change the **domain** layer under `internal/domain/`. This layer has no dependencies on other internal packages.

## Where to put things

- **Entities:** `internal/domain/entity/<name>.go` — one file per aggregate/entity (e.g. `user.go`, `appointment.go`).
- **Repository interfaces (ports):** `internal/domain/repository/<name>_repository.go` — interfaces only; infrastructure implements them.

## Entity checklist

- Struct with exported fields; use `entity` package.
- Allowed: stdlib, `time`, JSON/GORM struct tags for persistence.
- No imports from `application`, `infrastructure`, or `api`.
- Prefer `ID` as primary key (e.g. `uint` or `uuid.UUID`); include `CreatedAt`/`UpdatedAt` if needed.

**Example:**

```go
// internal/domain/entity/user.go
package entity

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

## Repository interface checklist

- Define the interface in `internal/domain/repository/`.
- Methods use `entity` types in parameters and returns (e.g. `*entity.User`, `[]entity.User`).
- Name: `<Entity>Repository` (e.g. `UserRepository`).
- Typical methods: `FindAll`, `FindByID`, `Create`, `Update`, `Delete` — add only what the use cases need.
- Add a short comment that the domain owns the interface and infrastructure implements it.

**Example:**

```go
// internal/domain/repository/user_repository.go
package repository

import "gin-quickstart/internal/domain/entity"

// UserRepository defines the contract for user data access.
// The domain owns this interface; infrastructure implements it.
type UserRepository interface {
	FindAll() ([]entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Create(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(id uint) error
}
```

## Do not

- Import `internal/application`, `internal/infrastructure`, or `internal/api`.
- Put repository **implementations** in domain; those go in `internal/infrastructure/repository/`.
- Add HTTP or framework-specific types; keep domain pure.

After domain changes, add or update the use case (application layer), then infrastructure (repository impl), then API and wiring in `main.go`.

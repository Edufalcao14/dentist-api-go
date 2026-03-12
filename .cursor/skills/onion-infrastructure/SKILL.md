---
name: onion-infrastructure
description: Implements the infrastructure layer in this project's onion architecture. Use when adding or changing database setup, repository implementations (adapters), or external integrations in internal/infrastructure/.
---

# Infrastructure Layer (Onion Architecture)

Implement or change the **infrastructure** layer under `internal/infrastructure/`. This layer provides concrete implementations of domain ports (e.g. repository implementations) and technical setup (DB connection).

## Where to put things

- **Database connection / migrations:** `internal/infrastructure/database/` (e.g. `postgres.go`).
- **Repository implementations (adapters):** `internal/infrastructure/repository/<entity>_repository.go` (e.g. `user_repository.go`).

## Repository implementation checklist

- Implement the **interface** defined in `internal/domain/repository/` (e.g. `repository.UserRepository`).
- Type name should reflect the technology (e.g. `GormUserRepository`).
- Constructor: `NewGorm<Entity>Repository(db *gorm.DB) *Gorm<Entity>Repository`.
- All method signatures must match the domain interface exactly (same parameter and return types using `entity` types).
- File may import: `domain/entity`, `domain/repository` (for interface), `gorm.io/gorm`, and project `config` if needed.

**Example:**

```go
// internal/infrastructure/repository/user_repository.go
package repository

import (
	"gin-quickstart/internal/domain/entity"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// ... Create, Update, Delete to match domain interface
```

## Database setup

- If adding a new entity that must be migrated, ensure the entity is registered in the same place where AutoMigrate is called (e.g. `database/postgres.go` or `main.go`), e.g. `db.AutoMigrate(&entity.User{}, &entity.Appointment{})`.
- Infrastructure may use `config` for connection settings; keep DB-specific code in `internal/infrastructure/`.

## Do not

- Define new “repository” interfaces in infrastructure; interfaces live in `domain/repository/`.
- Import `internal/application` or `internal/api`.
- Put business logic here; only persistence and technical implementation.

After infrastructure changes, wire the new or updated repository in `cmd/api/main.go` (create concrete repo, pass to use case constructor). Ensure the API layer has a handler that uses the use case.

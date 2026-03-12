---
name: onion-application
description: Implements the application (use-case) layer in this project's onion architecture. Use when adding or changing use cases, application DTOs, or business orchestration in internal/application/usecase/.
---

# Application Layer (Onion Architecture)

Implement or change the **application** layer under `internal/application/usecase/`. This layer orchestrates domain entities and repository ports only.

## Where to put things

- **Use cases:** `internal/application/usecase/<entity>_usecase.go` (e.g. `user_usecase.go`, `appointment_usecase.go`).
- **Input DTOs** for use case methods: define in the same file or same package (e.g. `CreateUserInput`, `UpdateUserInput`).
- **Domain errors** used by use case: define in usecase package (e.g. `ErrUserNotFound`, `ErrEmailAlreadyInUse`).

## Dependency rule

- Import **only** `internal/domain/entity` and `internal/domain/repository` (interfaces).
- Do **not** import `internal/infrastructure` or `internal/api`. No HTTP, no DB driver, no Gin.

## Use case struct and constructor

- Struct holds repository **interface** (e.g. `repo repository.UserRepository`).
- Constructor: `New<Entity>UseCase(repo repository.<Entity>Repository) *<Entity>UseCase`.
- Use the interface type from domain, not a concrete implementation.

**Example:**

```go
type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}
```

## Methods and DTOs

- One use case method per application operation (e.g. `GetAll`, `GetByID`, `Create`, `Update`, `Delete`).
- For create/update, define input structs (e.g. `CreateUserInput`, `UpdateUserInput`) with the fields the use case needs; handlers will map from HTTP request to these.
- Return domain entities (e.g. `*entity.User`, `[]entity.User`) or domain-relevant errors (e.g. `ErrUserNotFound`).
- Use case performs orchestration and validation that belongs to the application (e.g. “user must exist before update”); keep pure domain invariants in entity/domain if needed.

**Example:**

```go
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type CreateUserInput struct {
	Name  string
	Email string
}

func (uc *UserUseCase) Create(input CreateUserInput) (*entity.User, error) {
	user := &entity.User{Name: input.Name, Email: input.Email}
	return uc.repo.Create(user)
}
```

## Do not

- Import infrastructure or api; no `*gorm.DB`, no `gin.Context`.
- Return HTTP status or request/response DTOs; handlers will map use case results to HTTP.
- Put repository implementations here; use case depends only on the interface from `domain/repository`.

After application changes, ensure infrastructure has an implementation of the repository interface and that `main.go` wires repo → use case; then add or update the handler (API layer).

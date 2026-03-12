---
name: onion-api
description: Implements the API (presentation) layer in this project's onion architecture. Use when adding or changing HTTP handlers, routes, request/response DTOs, or middleware in internal/api/.
---

# API / Presentation Layer (Onion Architecture)

Implement or change the **API** layer under `internal/api/`. This layer exposes HTTP endpoints and delegates business logic to use cases.

## Where to put things

- **Handlers:** `internal/api/handler/<entity>_handler.go` (e.g. `user_handler.go`).
- **Request/response DTOs:** In the same handler file or package (e.g. `CreateUserRequest`, `UpdateUserRequest`) with JSON and `binding` tags for Gin.
- **Routes:** Register in `internal/api/router/router.go`.
- **Middleware:** `internal/api/middleware/` (e.g. error recovery, auth).

## Dependency rule

- Import **only** `internal/application/usecase` for business operations.
- Do **not** import `internal/domain` or `internal/infrastructure` in handlers for business flow (router may pass handlers that were constructed in main with use cases).

## Handler struct and constructor

- Handler holds the **use case** (e.g. `useCase *usecase.UserUseCase`).
- Constructor: `New<Entity>Handler(uc *usecase.<Entity>UseCase) *<Entity>Handler`.
- Handlers are created in `main.go` and passed to the router.

**Example:**

```go
type UserHandler struct {
	useCase *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}
```

## Request/response and mapping

- Define request structs for create/update with `json` and `binding` tags (e.g. `binding:"required"`, `binding:"required,email"`).
- Parse request in handler (e.g. `c.ShouldBindJSON(&req)`), map to use-case input (e.g. `usecase.CreateUserInput{Name: req.Name, Email: req.Email}`), call use case, then map result or error to HTTP status and JSON.
- Return domain entities as JSON when appropriate (they already have `json` tags); for errors, map use-case errors (e.g. `usecase.ErrUserNotFound`) to 404 and generic errors to 500.

**Example:**

```go
type CreateUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.useCase.Create(usecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
```

## Router

- Router receives handler(s) and registers routes (e.g. `GET/POST /api/v1/users`, `GET/PUT/DELETE /api/v1/users/:id`).
- Use the same URL prefix and versioning as existing routes; add Swagger comments if the project uses them.

## Do not

- Import domain or infrastructure in handlers for business logic; use only use case methods.
- Put business rules or repository access in handlers; keep that in the application layer.
- Skip the use case and call infrastructure from the handler.

After API changes, ensure `main.go` constructs the handler with the correct use case and that the router is given the new handler (or group). See `docs/ARCHITECTURE.md` for the full flow.

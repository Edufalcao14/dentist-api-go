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

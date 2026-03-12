package usecase

import (
	"errors"
	"gin-quickstart/internal/domain/entity"
	"gin-quickstart/internal/domain/repository"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	Name  string
	Email string
}

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetAll() ([]entity.User, error) {
	return uc.repo.FindAll()
}

func (uc *UserUseCase) GetByID(id uint) (*entity.User, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *UserUseCase) Create(input CreateUserInput) (*entity.User, error) {
	user := &entity.User{
		Name:  input.Name,
		Email: input.Email,
	}
	return uc.repo.Create(user)
}

func (uc *UserUseCase) Update(id uint, input UpdateUserInput) (*entity.User, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Name = input.Name
	user.Email = input.Email

	return uc.repo.Update(user)
}

func (uc *UserUseCase) Delete(id uint) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	return uc.repo.Delete(id)
}

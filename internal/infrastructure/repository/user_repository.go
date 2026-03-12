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

func (r *GormUserRepository) Create(user *entity.User) (*entity.User, error) {
	result := r.db.Create(user)
	return user, result.Error
}

func (r *GormUserRepository) Update(user *entity.User) (*entity.User, error) {
	result := r.db.Save(user)
	return user, result.Error
}

func (r *GormUserRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.User{}, id)
	return result.Error
}

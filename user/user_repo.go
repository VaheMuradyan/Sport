package user

import (
	"github.com/VaheMuradyan/Sport/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) CreateUser(user *models.User) error {
	result := repo.DB.Create(user)
	return result.Error
}

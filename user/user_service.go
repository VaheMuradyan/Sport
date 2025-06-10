package user

import (
	"github.com/VaheMuradyan/Sport/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (service *UserService) RegisterUser(req models.UserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return service.UserRepo.CreateUser(&user)
}

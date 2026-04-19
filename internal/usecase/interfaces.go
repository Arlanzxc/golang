package usecase

import (
	"GOLANG/internal/entity"
)

type (
	UserInterface interface {
		RegisterUser(user *entity.User) (*entity.User, string, error)
		LoginUser(user *entity.LoginUserDTO) (string, error)
		GetUserByID(id string) (*entity.User, error)
		PromoteUser(id string) error

		RegisterUserWithEmailCheck(user *entity.User, email string) error
		UpdateUserName(id string, newName string) error
		DeleteUser(id string) error
	}

	UserRepoInterface interface {
		RegisterUser(user *entity.User) (*entity.User, error)
		LoginUser(username string) (*entity.User, error)
		GetUserByID(id string) (*entity.User, error)
		PromoteUser(id string) error

		GetByEmail(email string) (*entity.User, error)
		UpdateUser(user *entity.User) error
		DeleteUser(id string) error
	}
)
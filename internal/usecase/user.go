package usecase

import (
	"GOLANG/internal/repository"
	"GOLANG/pkg/modules"
)

type UserUsecase interface {
	GetAll() ([]modules.User, error)
	GetByID(id int) (*modules.User, error)
	Create(u modules.User) (int, error)
	Update(u modules.User) error
	Delete(id int) (int64, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetAll() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *userUsecase) GetByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) Create(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *userUsecase) Update(user modules.User) error {
	return u.repo.UpdateUser(user)
}

func (u *userUsecase) Delete(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}

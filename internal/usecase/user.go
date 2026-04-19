package usecase

import (
	"fmt"

	"GOLANG/internal/entity"
	"GOLANG/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo UserRepoInterface
}

func NewUserUseCase(r UserRepoInterface) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = "user" 
	}

	createdUser, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}

	sessionID := uuid.New().String()
	return createdUser, sessionID, nil
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.LoginUser(input.Username)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if !utils.CheckPassword(userFromRepo.Password, input.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteUser(id)
}

func (u *UserUseCase) RegisterUserWithEmailCheck(user *entity.User, email string) error {
	existing, err := u.repo.GetByEmail(email)
	if existing != nil {
		return fmt.Errorf("user with this email already exists")
	}
	if err != nil && err.Error() != "user not found" {
		return fmt.Errorf("error getting user with this email")
	}
	_, err = u.repo.RegisterUser(user)
	return err
}

func (u *UserUseCase) UpdateUserName(id string, newName string) error {
	if newName == "" {
		return fmt.Errorf("name cannot be empty")
	}
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Username = newName 
	return u.repo.UpdateUser(user)
}

func (u *UserUseCase) DeleteUser(id string) error {
	if id == "1" { 
		return fmt.Errorf("it is not allowed to delete admin user")
	}
	return u.repo.DeleteUser(id)
}
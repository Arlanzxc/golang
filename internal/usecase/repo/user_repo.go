package repo

import (
	"fmt"

	"GOLANG/internal/entity"
	"GOLANG/pkg/postgres"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{PG: pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	err := u.PG.Conn.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(username string) (*entity.User, error) {
	var userFromDB entity.User
	if err := u.PG.Conn.Where("username = ?", username).First(&userFromDB).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) PromoteUser(id string) error {
	return u.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", "admin").Error
}

func (u *UserRepo) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (u *UserRepo) UpdateUser(user *entity.User) error {
	return u.PG.Conn.Save(user).Error
}

func (u *UserRepo) DeleteUser(id string) error {
	return u.PG.Conn.Where("id = ?", id).Delete(&entity.User{}).Error
}
package repository

import (
	"GOLANG/internal/repository/_postgres"
	"GOLANG/internal/repository/_postgres/users"
	"GOLANG/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(u modules.User) (int, error)
	UpdateUser(u modules.User) error
	DeleteUser(id int) (int64, error)
}
type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}

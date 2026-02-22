package users

import (
	"GOLANG/internal/repository/_postgres"
	"GOLANG/pkg/modules"
	"fmt"
	"time"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}
func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int
	query := "INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.DB.QueryRow(query, user.Name, user.Email, user.Age).Scan(&id)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var u modules.User
	err := r.db.DB.Get(&u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *Repository) UpdateUser(user modules.User) error {
	res, err := r.db.DB.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	res, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return rows, nil
}

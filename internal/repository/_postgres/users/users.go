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

func (r *Repository) GetPaginatedUsers(filters map[string]string, sortBy string, page, pageSize int) (modules.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	var users []modules.User
	var totalCount int

	baseQuery := " FROM users WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if val, ok := filters["id"]; ok && val != "" {
		baseQuery += fmt.Sprintf(" AND id = $%d", argID)
		args = append(args, val)
		argID++
	}
	if val, ok := filters["name"]; ok && val != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+val+"%")
		argID++
	}
	if val, ok := filters["email"]; ok && val != "" {
		baseQuery += fmt.Sprintf(" AND email ILIKE $%d", argID)
		args = append(args, "%"+val+"%")
		argID++
	}
	if val, ok := filters["gender"]; ok && val != "" {
		baseQuery += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, val)
		argID++
	}
	if val, ok := filters["birth_date"]; ok && val != "" {
		baseQuery += fmt.Sprintf(" AND birth_date = $%d", argID)
		args = append(args, val)
		argID++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	r.db.DB.Get(&totalCount, countQuery, args...)

	orderQuery := " ORDER BY id ASC"
	if sortBy != "" {
		orderQuery = fmt.Sprintf(" ORDER BY %s ASC", sortBy)
		if sortBy[0] == '-' {
			orderQuery = fmt.Sprintf(" ORDER BY %s DESC", sortBy[1:])
		}
	}

	selectQuery := "SELECT * " + baseQuery + orderQuery + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, pageSize, offset)

	err := r.db.DB.Select(&users, selectQuery, args...)
	if err != nil {
		return modules.PaginatedResponse{}, err
	}

	if users == nil {
		users = []modules.User{}
	}

	return modules.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]modules.User, error) {
	var friends []modules.User
	query := `
		SELECT u.* FROM users u
		JOIN user_friends uf1 ON u.id = uf1.friend_id
		JOIN user_friends uf2 ON u.id = uf2.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
	`
	err := r.db.DB.Select(&friends, query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	if friends == nil {
		friends = []modules.User{}
	}
	return friends, nil
}
package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"` 
	Role     string    `json:"role"`
	Verified bool      `json:"verified"`
}
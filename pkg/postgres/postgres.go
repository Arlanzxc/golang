package postgres

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Conn *gorm.DB
}

func NewPostgres(dsn string) *Postgres {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &Postgres{Conn: db}
}
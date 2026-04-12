	package app

import (
	"fmt"
	"log"
	"os"

	v1 "GOLANG/internal/controller/http/v1"
	"GOLANG/internal/entity"
	"GOLANG/internal/usecase"
	"GOLANG/internal/usecase/repo"
	"GOLANG/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5252"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres159357"),
		getEnv("DB_NAME", "mydb"),
		getEnv("DB_SSLMODE", "disable"),
	)

	pg := postgres.NewPostgres(dsn)
	fmt.Println("Connected to PostgreSQL via GORM")

	pg.Conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.Default()

	apiV1 := router.Group("/v1")
	v1.NewUserRoutes(apiV1, userUseCase)

	log.Println("Server starting on :8090")
	if err := router.Run(":8090"); err != nil { 
		log.Fatalf("Server failed: %v", err)
	}
}
package app

import (
	"GOLANG/internal/handler"
	"GOLANG/internal/repository/_postgres"
	"GOLANG/internal/repository/_postgres/users"
	"GOLANG/internal/usecase"
	"GOLANG/pkg/modules"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()

	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	fmt.Println("Connected to database:", _postgre)

	repo := users.NewUserRepository(_postgre)
	uc := usecase.NewUserUsecase(repo)
	h := handler.NewUserHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/users", h.HandleUsers)
	mux.HandleFunc("/users/", h.HandleUsers)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", middleware(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		if r.Header.Get("X-API-KEY") != "KBTU_SECRET_2026" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5252"),
		Username:    getEnv("DB_USER", "postgres"),
		Password:    getEnv("DB_PASSWORD", "postgres159357"),
		DBName:      getEnv("DB_NAME", "mydb"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}

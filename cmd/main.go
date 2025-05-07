package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juancanchi/users/internal/delivery/http"
	"github.com/juancanchi/users/internal/infrastructure/postgres"
	"github.com/juancanchi/users/internal/usecase"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=jujuy_market port=5432 sslmode=disable"
	db, err := postgres.NewDB(dsn)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	repo := postgres.NewUserRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecreto" // para dev
	}

	uc := usecase.NewUserUsecase(repo, jwtSecret)
	handler := http.NewUserHandler(uc)

	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🧑‍💻 User service listening on http://localhost:%s", port)
	r.Run(":" + port)
}

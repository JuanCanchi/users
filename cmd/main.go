package main

import (
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juancanchi/users/internal/delivery/http"
	"github.com/juancanchi/users/internal/infrastructure/postgres"
	"github.com/juancanchi/users/internal/usecase"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=jujuy_market port=5432 sslmode=disable"
	}

	var dbErr error
	var db = tryConnectToDB(dsn, &dbErr)
	if dbErr != nil {
		log.Fatalf("❌ Error al conectar a la base de datos después de varios intentos: %v", dbErr)
	}

	repo := postgres.NewUserRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecreto" // para dev
	}

	uc := usecase.NewUserUsecase(repo, jwtSecret)
	handler := http.NewUserHandler(uc)

	r := gin.Default()

	// 🛡️ Habilitar CORS para permitir conexión desde el frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🧑‍💻 User service listening on http://localhost:%s", port)
	r.Run(":" + port)
}

func tryConnectToDB(dsn string, lastErr *error) *gorm.DB {
	const maxAttempts = 10
	for i := 1; i <= maxAttempts; i++ {
		log.Printf("⏳ Intentando conectar a la base de datos... intento %d/%d", i, maxAttempts)
		db, err := postgres.NewDB(dsn)
		if err == nil {
			log.Println("✅ Conexión a la base de datos exitosa.")
			return db
		}
		*lastErr = err
		time.Sleep(3 * time.Second)
	}
	return nil
}

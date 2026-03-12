// @title           Dentist API
// @version         1.0
// @description     Dentist API built with Gin following Onion Architecture.

// @host      localhost:8080
// @BasePath  /api/v1

package main

import (
	"log"

	"gin-quickstart/config"
	"gin-quickstart/internal/api/handler"
	"gin-quickstart/internal/api/router"
	"gin-quickstart/internal/application/usecase"
	"gin-quickstart/internal/infrastructure/database"
	infrarepo "gin-quickstart/internal/infrastructure/repository"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := infrarepo.NewGormUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUC)

	r := router.New(userHandler)

	log.Printf("server running on port %s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

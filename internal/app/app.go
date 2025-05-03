package app

import (
	"fmt"
	"log"
	"os"
	"routinist/internal/domain/model"
	"routinist/internal/seed"

	"routinist/internal/controller/http"
	"routinist/internal/repository"
	"routinist/internal/usecase"
	"routinist/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	// Initialize logger
	l := logger.New("app")

	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to database
	dbpool, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = dbpool.AutoMigrate(&model.User{}, &model.Unit{}, &model.Habit{}, &model.HabitUnit{}, &model.UserHabit{})
	if err != nil {
		log.Fatalf("Failed to migrations database: %v", err)
	}

	seed.Seed(dbpool, l)

	// Initialize Gin router
	router := gin.Default()
	authRepo := repository.NewAuthRepo(dbpool, l)
	habitRepo := repository.NewHabitRepo(dbpool, l)

	// Initialize usecase
	authUseCase := usecase.NewAuthUseCase(authRepo, l)
	habitUseCase := usecase.NewHabitUseCase(habitRepo, l)

	// Setup routes
	http.NewRouter(router, l, authUseCase, habitUseCase)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

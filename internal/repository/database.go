package repository

import (
	"log"
	"os"
	"routinist/internal/domain/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
}

func AutoMigrate() {
	err := DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

package config

import (
	"os"

	"reverse-engineering-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/reverse_engineering_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自動マイグレーション
	err = db.AutoMigrate(
		&models.Project{},
		&models.File{},
		&models.Analysis{},
		&models.User{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

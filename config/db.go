package config

import (
	"fmt"
	"log"
	"os"
	"rental-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "testdb"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Successfully connected to the database.")

	if err := autoMigrateModels(); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	log.Println("Database migrated successfully.")
	return DB
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func autoMigrateModels() error {
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.MesinBor{},
		&models.RentalHistory{},
		&models.Review{},
		&models.Maintenance{},
	}

	if err := DB.AutoMigrate(modelsToMigrate...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}
	return nil
}

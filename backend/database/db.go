package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// InitDB initializes the connection to PostgreSQL and runs auto-migrations
func InitDB(dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic sql.DB: %w", err)
	}

	// Connection pool tuning
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	log.Println("Database connection pool established successfully.")

	// Automatically run migrations
	err = db.AutoMigrate(&SystemConfig{}, &Token{}, &Trade{}, &User{})
	if err != nil {
		return nil, fmt.Errorf("database auto-migration failed: %w", err)
	}

	log.Println("Database auto-migrations executed successfully.")
	DB = db
	return db, nil
}

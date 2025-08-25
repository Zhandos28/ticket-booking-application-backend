package db

import (
	"fmt"

	"github.com/Zhandos28/ticket-booking/config"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(config *config.EnvConfig, DBMigrator func(db *gorm.DB) error) *gorm.DB {
	uri := fmt.Sprintf(
		`host=%s user=%s dbname=%s sslmode=%s password=%s port=5432`,
		config.DBHost, config.DBUser, config.DBName, config.DBSSLMode, config.DBPassword,
	)

	fmt.Println(uri)
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Unable to connect to database: %e", err)
	}

	log.Info("Successfully connected to database")

	if err := DBMigrator(db); err != nil {
		log.Fatalf("Unable to migrate tables: %e", err)
	}

	return db
}

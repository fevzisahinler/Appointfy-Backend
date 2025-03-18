package db

import (
	"context"
	"fmt"
	"github.com/fevzisahinler/Appointfy-Backend/models"
	"time"

	"github.com/fevzisahinler/Appointfy-Backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGDBName, cfg.PGPort,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Failed to get database instance:", err)
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.PingContext(ctx); err != nil {
		fmt.Println("Database ping failed:", err)
		return err
	}
	if err := AutoMigrate(db); err != nil {
		return err
	}

	fmt.Println("Database connected successfully")

	DB = db
	return nil
}

func AutoMigrate(database *gorm.DB) error {
	if err := database.AutoMigrate(
		&models.User{},
		&models.Role{},
	); err != nil {
		fmt.Println("Failed to migrate models")
		return err
	}
	return nil
}

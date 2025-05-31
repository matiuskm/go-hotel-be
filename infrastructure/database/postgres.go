package database

import (
	"log"
	"matiuskm/go-hotel-be/domain/entities"
	"matiuskm/go-hotel-be/pkg/seeders"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entities.User{}); err!= nil {
		return nil, err
	}
	seeders.SeedAdmin(db)
	log.Println("Database migrated and admin seeded")
	return db, nil
}
package seeders

import (
	"log"
	"matiuskm/go-hotel-be/domain/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	adminUsername := "" // your admin username here
	adminPassword := "" // your admin password here

	var admin entities.User
	result := db.First(&admin, "username = ?", adminUsername)
	if result.Error == nil {
		log.Println("Admin user already exists")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	admin = entities.User{
		Username: adminUsername,
		Password: string(hashed),
		FullName: "Lazy Coder",
		Role: "admin",
	}
	db.Create(&admin)
	log.Println("Admin user created")
}
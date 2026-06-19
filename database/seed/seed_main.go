package seed

import (
	"log"

	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/utils"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin123")
	admin := models.User{
		Name: "Super Admin",
		Email: "admin@example.com",
		Password: password,
		Role: "Admin",
	}

	if err := config.DB.FirstOrCreate(&admin,models.User{Email:admin.Email}).Error; err != nil {
		log.Println("Failed to seed admin", err)
	} else {
		log.Println("Admin user seeded")
	}
}
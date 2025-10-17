package config

import (
	"log"

	"go-api/models"
)

func InitMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Session{},
	)
	if err != nil {
		log.Fatalf("❌ Gagal melakukan migrasi: %v", err)
	}
	log.Println("✅ Migrasi database selesai.")
}

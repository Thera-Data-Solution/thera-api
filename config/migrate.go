package config

import (
	"log"
)

func InitMigrate() {
	// err := DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Session{},
	// )
	// if err != nil {
	// 	log.Fatalf("❌ Gagal melakukan migrasi: %v", err)
	// }
	log.Println("✅ Migrasi database tidak dilakukan, mengingat service nextjs sedang berjalan.")
}

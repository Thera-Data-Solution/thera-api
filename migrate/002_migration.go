package migrate

import (
	"log"
	"thera-api/config"
	"thera-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration002() {
	db := config.DB
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "003_add_password_reset_token",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&models.PasswordResetToken{},
				)
			},
			Rollback: nil,
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("❌ Initial migration failed:", err)
	}
	log.Println("✅ Initial migration completed")
}

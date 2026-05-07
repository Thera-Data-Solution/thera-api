package migrate

import (
	"log"
	"thera-api/config"
	"thera-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration004() {
	db := config.DB
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "004_adding_disable_column",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&models.Categories{},
					&models.Booked{},
					&models.Review{},
					&models.Events{},
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

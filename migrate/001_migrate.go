package migrate

import (
	"log"
	"thera-api/config"
	"thera-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration001() {
	db := config.DB
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "001_init_all_models",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&models.Article{},
					&models.Booked{},
					&models.Categories{},
					&models.Gallery{},
					&models.Hero{},
					&models.Link{},
					&models.Partner{},
					&models.PositionLanding{},
					&models.ResetPasswordRequest{},
					&models.Review{},
					&models.Schedules{},
					&models.Session{},
					&models.Setting{},
					&models.TenantUser{},
					&models.Tenant{},
					&models.Translation{},
					&models.User{},
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

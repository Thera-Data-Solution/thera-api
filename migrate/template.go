package migrate

import (
	"log"
	"thera-api/config"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// TemplateMigration digunakan sebagai contoh untuk menambahkan field baru ke model apapun
// Cara pakai:
// 1. Copy file ini dan ganti nama file (misal 003_add_xxx_field.go)
// 2. Ganti ID migrasi unik
// 3. Masukkan model yang diubah di AutoMigrate
// 4. Jalankan Rollback jika ingin menghapus kolom
func TemplateMigration(id string, model interface{}, rollbackColumn string) {
	db := config.DB

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: id,
			Migrate: func(tx *gorm.DB) error {
				// AutoMigrate hanya untuk model yang diubah
				return tx.AutoMigrate(model)
			},
			Rollback: func(tx *gorm.DB) error {
				if rollbackColumn != "" {
					// Drop kolom jika rollback
					return tx.Migrator().DropColumn(model, rollbackColumn)
				}
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("❌ Migration %s failed: %v", id, err)
	}

	log.Printf("✅ Migration %s completed\n", id)
}

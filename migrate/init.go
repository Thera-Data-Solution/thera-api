package migrate

import "thera-api/models"

func RunMigrations() {
	// jalankan migrasi urut berdasarkan ID
	Migration001() // dari 001_init.go
	TemplateMigration("002_adding_preload", &models.Categories{}, "")
	TemplateMigration("002_adding_preload_schedules", &models.Schedules{}, "")
	TemplateMigration("003_adding_preload_bookings", &models.Booked{}, "")

	// dst.
}

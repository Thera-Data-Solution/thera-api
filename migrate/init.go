package migrate

import "thera-api/models"

func RunMigrations() {
	Migration001()
	TemplateMigration("002_adding_preload", &models.Categories{}, "")
	TemplateMigration("002_adding_preload_schedules", &models.Schedules{}, "")
	TemplateMigration("003_adding_preload_bookings", &models.Booked{}, "")
	TemplateMigration("004_adding_custom_field", &models.Categories{}, "")
	TemplateMigration("005_adding_custom_field", &models.CategoryCustomField{}, "")
	TemplateMigration("006_adding_custom_field", &models.Booked{}, "")
	TemplateMigration("007_adding_custom_field", &models.BookedCustomField{}, "")
}

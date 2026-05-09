package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type EventsRepository struct {
	DB *gorm.DB
}

func (r *EventsRepository) Create(event *models.Events) error {
	return r.DB.Create(event).Error
}

func (r *EventsRepository) FindAll(tenant string) ([]models.Events, error) {
	var events []models.Events
	err := r.DB.
		Where(`tenant_id = ? AND status = ?`, tenant, "available").
		Order("start_at ASC").
		Find(&events).
		Error
	return events, err
}

func (r *EventsRepository) FindAllAsAdmin(tenant string) ([]models.Events, error) {
	var events []models.Events
	err := r.DB.
		Where(`tenant_id = ?`, tenant).
		Order("start_at ASC").
		Find(&events).
		Error
	return events, err
}

func (r *EventsRepository) FindByID(id string, tenant string) (*models.Events, error) {
	var event models.Events
	err := r.DB.First(&event, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventsRepository) FindBySlug(slug string) (*models.Events, error) {
	var event models.Events
	err := r.DB.First(&event, `slug = ?`, slug).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventsRepository) Update(event *models.Events) error {
	return r.DB.Save(event).Error
}

func (r *EventsRepository) Delete(id string, tenantId string) error {
	return r.DB.Delete(&models.Events{}, `id = ? AND tenant_id = ?`, id, tenantId).Error
}

package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventsService struct {
	EventsRepo *repositories.EventsRepository
}

func (s *EventsService) GetAllEvents(tenant string) ([]models.Events, error) {
	return s.EventsRepo.FindAll(tenant)
}

func (s *EventsService) GetAllEventsAsAdmin(tenant string) ([]models.Events, error) {
	return s.EventsRepo.FindAllAsAdmin(tenant)
}

func (s *EventsService) GetEventByID(id string, tenant string) (*models.Events, error) {
	return s.EventsRepo.FindByID(id, tenant)
}

func (s *EventsService) CreateEvent(
	name, nameEn string,
	description, descriptionEn *string,
	slug string,
	image *string, // Tambah ini
	price float64,
	startAt, endAt time.Time,
	capacity int,
	status string,
	tenantId *string,
	customFields datatypes.JSON,
	location string,
) (*models.Events, error) {
	if tenantId == nil || strings.TrimSpace(*tenantId) == "" {
		return nil, errors.New("tenantId tidak boleh kosong")
	}

	finalSlug, err := s.ensureUniqueSlug(s.pickSlug(slug, name), nil)
	if err != nil {
		return nil, err
	}

	event := &models.Events{
		Name:          name,
		NameEn:        nameEn,
		Description:   description,
		DescriptionEn: descriptionEn,
		Image:         image, // Tambah ini
		Slug:          finalSlug,
		Price:         price,
		StartAt:       startAt,
		EndAt:         endAt,
		Capacity:      capacity,
		Status:        status,
		TenantId:      tenantId,
		CustomFields:  customFields,
		Location:      location,
	}

	if err := s.EventsRepo.Create(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventsService) UpdateEvent(
	id string,
	tenantId string,
	updates map[string]interface{},
) (*models.Events, error) {
	event, err := s.EventsRepo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}

	if val, ok := updates["name"].(string); ok {
		event.Name = val
		newSlug, _ := s.ensureUniqueSlug(s.slugify(val), &event.ID)
		event.Slug = newSlug
	}
	if val, ok := updates["nameEn"].(string); ok {
		event.NameEn = val
	}
	if val, ok := updates["description"].(*string); ok {
		event.Description = val
	}
	if val, ok := updates["descriptionEn"].(*string); ok {
		event.DescriptionEn = val
	}
	if val, ok := updates["image"].(*string); ok {
		event.Image = val
	} // Tambah ini
	if val, ok := updates["price"].(float64); ok {
		event.Price = val
	}
	if val, ok := updates["startAt"].(time.Time); ok {
		event.StartAt = val
	}
	if val, ok := updates["endAt"].(time.Time); ok {
		event.EndAt = val
	}
	if val, ok := updates["capacity"].(int); ok {
		event.Capacity = val
	}
	if val, ok := updates["status"].(string); ok {
		event.Status = val
	}
	if val, ok := updates["customFields"].(datatypes.JSON); ok {
		event.CustomFields = val
	}

	if err := s.EventsRepo.Update(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventsService) DeleteEvent(id string, tenantId string) error {
	return s.EventsRepo.Delete(id, tenantId)
}

func (s *EventsService) pickSlug(slug, fallback string) string {
	if trimmed := strings.TrimSpace(slug); trimmed != "" {
		return trimmed
	}
	return fallback
}

func (s *EventsService) slugify(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	slug := re.ReplaceAllString(lower, "-")
	return strings.Trim(slug, "-")
}

func (s *EventsService) ensureUniqueSlug(base string, excludeID *string) (string, error) {
	seedSlug := s.slugify(base)
	slugCandidate := seedSlug
	addedTimestamp := false

	for {
		existing, err := s.EventsRepo.FindBySlug(slugCandidate)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		if existing == nil || (excludeID != nil && existing.ID == *excludeID) {
			return slugCandidate, nil
		}
		if !addedTimestamp {
			slugCandidate = fmt.Sprintf("%s-%d", seedSlug, time.Now().Unix())
			addedTimestamp = true
			continue
		}
		slugCandidate = fmt.Sprintf("%s-%d", seedSlug, time.Now().UnixNano())
	}
}

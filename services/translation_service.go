package services

import (
	"math"
	"thera-api/models"
	"thera-api/repositories"
)

type TranslationService struct {
	Repo *repositories.TranslationRepository
}

// GetAllTranslations dengan advanced filter + pagination
func (s *TranslationService) GetAllTranslations(tenantId string, filter models.TranslationFilter, page, limit int) (*models.PaginatedTranslations, error) {
	// normalisasi page & limit di level service supaya konsisten
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	// batasi maksimal, tapi masih cukup besar untuk use case real
	if limit > 1000 {
		limit = 1000
	}

	items, total, err := s.Repo.FindFiltered(tenantId, filter, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &models.PaginatedTranslations{
		Data:       items,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s *TranslationService) GetTranslationByID(id string, tenantId string) (*models.Translation, error) {
	return s.Repo.FindByID(id, tenantId)
}

func (s *TranslationService) CreateTranslation(locale, namespace, key, value, tenantId string) (*models.Translation, error) {
	t := &models.Translation{
		Locale:    locale,
		Namespace: namespace,
		Key:       key,
		Value:     value,
		TenantId:  &tenantId,
	}
	if err := s.Repo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TranslationService) UpdateTranslation(id string, locale, namespace, key, value *string, tenantId string) (*models.Translation, error) {
	t, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}

	if locale != nil {
		t.Locale = *locale
	}
	if namespace != nil {
		t.Namespace = *namespace
	}
	if key != nil {
		t.Key = *key
	}
	if value != nil {
		t.Value = *value
	}

	if err := s.Repo.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TranslationService) DeleteTranslation(id string, tenantId string) error {
	return s.Repo.Delete(id, tenantId)
}

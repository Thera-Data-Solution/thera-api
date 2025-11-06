package services

import (
	"thera-api/models"
	"thera-api/repositories"
)

type TranslationService struct {
	Repo *repositories.TranslationRepository
}

func (s *TranslationService) GetAllTranslations(tenantId string) ([]models.Translation, error) {
	return s.Repo.FindAll(tenantId)
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

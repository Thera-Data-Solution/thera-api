package services

import (
	"thera-api/dto"
	"thera-api/models"
	"thera-api/repositories"

	"github.com/mashingan/smapping"
)

type SettingService interface {
	FindAll() ([]models.Setting, error)
	FindById(id string) (*models.Setting, error)
	FindByTenantId(tenantId string) (*models.Setting, error)
	Upsert(dto dto.SettingRequestBody) (*models.Setting, error)
	Delete(id string) error
}

type settingService struct {
	repo repositories.SettingRepo
}

func NewSettingService(repo repositories.SettingRepo) SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) FindAll() ([]models.Setting, error) {
	return s.repo.FindAll()
}

func (s *settingService) FindById(id string) (*models.Setting, error) {
	return s.repo.FindById(id)
}

func (s *settingService) FindByTenantId(tenantId string) (*models.Setting, error) {
	return s.repo.FindByTenantId(tenantId)
}

func (s *settingService) Upsert(dto dto.SettingRequestBody) (*models.Setting, error) {
	setting := models.Setting{}
	if err := smapping.FillStruct(&setting, smapping.MapFields(&dto)); err != nil {
		return nil, err
	}
	return s.repo.Upsert(&setting)
}

func (s *settingService) Delete(id string) error {
	return s.repo.Delete(id)
}
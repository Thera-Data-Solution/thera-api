package services

import (
	"thera-api/models"
	"thera-api/repositories"
	"time"
)

type TenantService struct {
	TenantRepo *repositories.TenantRepository
}

func (s *TenantService) GetAllTenants() ([]models.Tenant, error) {
	tenants, err := s.TenantRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s *TenantService) CreateTenant(name string, logo *string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		Name:      name,
		Logo:      logo,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.TenantRepo.Create(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *TenantService) GetTenantByID(id string) (*models.Tenant, error) {
	tenant, err := s.TenantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *TenantService) UpdateTenant(id string, name, logo *string, isActive *bool) (*models.Tenant, error) {

	tenant, err := s.TenantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		tenant.Name = *name
	}
	if logo != nil {
		tenant.Logo = logo
	}
	if isActive != nil {
		tenant.IsActive = *isActive
	}
	tenant.UpdatedAt = time.Now()

	if err := s.TenantRepo.Update(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *TenantService) DeleteTenant(id string) error {
	return s.TenantRepo.Delete(id)
}

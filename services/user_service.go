package services

import (
	"thera-api/models"
	"thera-api/repositories"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func (s *UserService) GetAllByTenantId(tenantId string) ([]models.User, error) {
	users, err := s.UserRepo.FindAllByTenantId(tenantId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetAllByTenantPaginated(tenantId string, page, pageSize int) ([]models.User, int64, error) {
	return s.UserRepo.GetAllByTenantPaginated(tenantId, page, pageSize)
}

func (s *UserService) GetAllPaginated(page, pageSize int) ([]models.User, []string, int64, error) {
	return s.UserRepo.GetAllPaginated(page, pageSize)
}

func (s *UserService) DisableUser(id, tenantId string) error {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return err
	}
	user.Disable = !user.Disable
	return s.UserRepo.Update(user)
}

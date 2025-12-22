package services

import (
	"errors"
	"thera-api/models"
	"thera-api/repositories"
)

type HeroService struct {
	Repo *repositories.HeroRepository
}

func (s *HeroService) CreateHero(
	title string,
	subtitle,
	description,
	imageURL *string,
	buttonText, buttonLink, themeType *string,
	isActive bool,
	tenantId *string,
) (*models.Hero, error) {

	if title == "" {
		return nil, errors.New("judul wajib diisi")
	}

	hero := &models.Hero{
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       imageURL,
		ButtonText:  buttonText,
		ButtonLink:  buttonLink,
		ThemeType:   themeType,
		IsActive:    isActive,
		TenantId:    tenantId,
	}

	err := s.Repo.Create(tenantId, hero)
	return hero, err
}

func (s *HeroService) GetAllHeroes(tenantId string) (models.Hero, error) {
	return s.Repo.FindAll(tenantId)
}

func (s *HeroService) GetHeroByID(id string, tenantId string) (*models.Hero, error) {
	return s.Repo.FindByID(id, tenantId)
}

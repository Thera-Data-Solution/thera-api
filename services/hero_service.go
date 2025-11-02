package services

import (
	"errors"
	"thera-api/models"
	"thera-api/repositories"

	"github.com/google/uuid"
)

type HeroService struct {
	Repo *repositories.HeroRepository
}

func (s *HeroService) CreateHero(
	title string,
	subtitle, description, image, buttonText, buttonLink, themeType *string,
	isActive bool,
	tenantId *string,
) (*models.Hero, error) {

	if title == "" {
		return nil, errors.New("judul wajib diisi")
	}

	hero := &models.Hero{
		ID:          uuid.NewString(),
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       image,
		ButtonText:  buttonText,
		ButtonLink:  buttonLink,
		ThemeType:   themeType,
		IsActive:    isActive,
		TenantId:    tenantId,
	}

	err := s.Repo.Create(hero)
	return hero, err
}

func (s *HeroService) GetAllHeroes(tenantId string) ([]models.Hero, error) {
	return s.Repo.FindAll(tenantId)
}

func (s *HeroService) GetHeroByID(id string, tenantId string) (*models.Hero, error) {
	return s.Repo.FindByID(id, tenantId)
}

func (s *HeroService) UpdateHero(
	id string,
	title *string,
	subtitle, description, image, buttonText, buttonLink, themeType *string,
	isActive *bool,
	tenantId string,
) (*models.Hero, error) {

	hero, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		return nil, errors.New("hero tidak ditemukan")
	}

	if title != nil {
		hero.Title = *title
	}
	if subtitle != nil {
		hero.Subtitle = subtitle
	}
	if description != nil {
		hero.Description = description
	}
	if image != nil {
		hero.Image = image
	}
	if buttonText != nil {
		hero.ButtonText = buttonText
	}
	if buttonLink != nil {
		hero.ButtonLink = buttonLink
	}
	if themeType != nil {
		hero.ThemeType = themeType
	}
	if isActive != nil {
		hero.IsActive = *isActive
	}

	err = s.Repo.Update(hero)
	return hero, err
}

func (s *HeroService) DeleteHero(id string, tenantId string) error {
	return s.Repo.Delete(id, tenantId)
}

package services

import (
	"thera-api/models"
	"thera-api/repositories"
)

type CategoriesService struct {
	CategoriesRepo *repositories.CategoriesRepository
}

func (s *CategoriesService) GetAllCategories(tenant string) ([]models.Categories, error) {
	return s.CategoriesRepo.FindAll(tenant)
}

func (s *CategoriesService) GetCategoryByID(id string, tenant string) (*models.Categories, error) {
	return s.CategoriesRepo.FindByID(id, tenant)
}

func (s *CategoriesService) GetCategoryByIDAndTenant(id string, tenant string) (*models.Categories, error) {
	return s.CategoriesRepo.FindByIDAndTenant(id, tenant)
}

func (s *CategoriesService) CreateCategory(
	name string,
	description, descriptionEn *string,
	slug string,
	image *string,
	start, end int,
	location *string,
	price *float64,
	isGroup, isFree, isPayAsYouWish, isManual, disable bool,
	tenantId *string,
) (*models.Categories, error) {
	category := &models.Categories{
		Name:           name,
		Description:    description,
		DescriptionEn:  descriptionEn,
		Slug:           slug,
		Image:          image,
		Start:          start,
		End:            end,
		Location:       location,
		Price:          price,
		IsGroup:        isGroup,
		IsFree:         isFree,
		IsPayAsYouWish: isPayAsYouWish,
		IsManual:       isManual,
		Disable:        disable,
		TenantId:       tenantId,
	}
	if err := s.CategoriesRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoriesService) UpdateCategory(
	id string,
	name, description, descriptionEn, slug, image *string,
	start, end *int,
	location *string,
	price *float64,
	isGroup, isFree, isPayAsYouWish, isManual, disable *bool,
	tenantId string,
) (*models.Categories, error) {
	category, err := s.CategoriesRepo.FindByID(id, tenantId)

	if err != nil {
		return nil, err
	}

	if name != nil {
		category.Name = *name
	}
	if description != nil {
		category.Description = description
	}
	if descriptionEn != nil {
		category.DescriptionEn = descriptionEn
	}
	if slug != nil {
		category.Slug = *slug
	}
	if image != nil {
		category.Image = image
	}
	if start != nil {
		category.Start = *start
	}
	if end != nil {
		category.End = *end
	}
	if location != nil {
		category.Location = location
	}
	if price != nil {
		category.Price = price
	}
	if isGroup != nil {
		category.IsGroup = *isGroup
	}
	if isFree != nil {
		category.IsFree = *isFree
	}
	if isPayAsYouWish != nil {
		category.IsPayAsYouWish = *isPayAsYouWish
	}
	if isManual != nil {
		category.IsManual = *isManual
	}
	if disable != nil {
		category.Disable = *disable
	}

	if err := s.CategoriesRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoriesService) DeleteCategory(id string, tenantId string) error {
	return s.CategoriesRepo.Delete(id, tenantId)
}

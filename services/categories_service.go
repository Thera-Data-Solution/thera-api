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

type CategoriesService struct {
	CategoriesRepo *repositories.CategoriesRepository
}

func (s *CategoriesService) GetAllCategories(tenant string) ([]models.Categories, error) {
	return s.CategoriesRepo.FindAll(tenant)
}

func (s *CategoriesService) GetAllCategoriesWithType(tenant string, id string) ([]models.Categories, error) {
	return s.CategoriesRepo.FindAllByType(tenant, id)
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
	customFields datatypes.JSON,
) (*models.Categories, error) {
	if tenantId == nil || strings.TrimSpace(*tenantId) == "" {
		return nil, errors.New("tenantId tidak boleh kosong")
	}

	finalSlug, err := s.ensureUniqueSlug(s.pickSlug(slug, name), nil)
	if err != nil {
		return nil, err
	}

	category := &models.Categories{
		Name:           name,
		Description:    description,
		DescriptionEn:  descriptionEn,
		Slug:           finalSlug,
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
		CustomFields:   customFields,
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
	customFields datatypes.JSON,
) (*models.Categories, error) {
	category, err := s.CategoriesRepo.FindByID(id, tenantId)

	if err != nil {
		return nil, err
	}

	if customFields != nil {
		category.CustomFields = customFields
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

	slugSource := category.Slug
	if slug != nil && strings.TrimSpace(*slug) != "" {
		slugSource = *slug
	} else if name != nil && strings.TrimSpace(*name) != "" {
		slugSource = *name
	}

	finalSlug, err := s.ensureUniqueSlug(s.slugify(slugSource), &category.ID)
	if err != nil {
		return nil, err
	}
	category.Slug = finalSlug

	if err := s.CategoriesRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoriesService) DeleteCategory(id string, tenantId string) error {
	return s.CategoriesRepo.Delete(id, tenantId)
}

func (s *CategoriesService) pickSlug(slug, fallback string) string {
	if trimmed := strings.TrimSpace(slug); trimmed != "" {
		return trimmed
	}
	return fallback
}

func (s *CategoriesService) slugify(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	slug := re.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "kategori"
	}
	return slug
}

func (s *CategoriesService) ensureUniqueSlug(base string, excludeID *string) (string, error) {
	seedSlug := s.slugify(base)
	baseSlug := seedSlug
	slugCandidate := baseSlug
	counter := 1
	addedTimestamp := false

	for {
		existing, err := s.CategoriesRepo.FindBySlug(slugCandidate)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		if existing == nil || (excludeID != nil && existing.ID == *excludeID) {
			return slugCandidate, nil
		}

		if !addedTimestamp {
			baseSlug = fmt.Sprintf("%s-%d", seedSlug, time.Now().Unix())
			slugCandidate = baseSlug
			addedTimestamp = true
			continue
		}

		slugCandidate = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
}

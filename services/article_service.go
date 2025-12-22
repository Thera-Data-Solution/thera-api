package services

import (
	"errors"
	"math"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct {
	Repo *repositories.ArticleRepository
}

// NewArticleService creates a new instance of ArticleService
func NewArticleService(repo *repositories.ArticleRepository) *ArticleService {
	return &ArticleService{Repo: repo}
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(
	title, slug string,
	coverImage, body, excerpt, articleType *string,
	isPublished bool,
	tenantID *string,
) (*models.Article, error) {
	if tenantID == nil {
		return nil, errors.New("tenantId tidak boleh kosong")
	}

	// Check if slug already exists for this tenant
	existingArticle, err := s.Repo.FindBySlug(slug, *tenantID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Failed to check existing slug", zap.String("slug", slug), zap.String("tenantId", *tenantID), zap.Error(err))
		return nil, err
	}
	if existingArticle != nil {
		return nil, errors.New("slug sudah ada untuk tenant ini")
	}

	article := &models.Article{
		Title:       title,
		Slug:        slug,
		CoverImage:  coverImage,
		Body:        body,
		Excerpt:     excerpt,
		IsPublished: isPublished,
		ArticleType: articleType,
		TenantId:    tenantID,
	}

	if err := s.Repo.Create(article); err != nil {
		logger.Log.Error("Failed to create article", zap.String("title", title), zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Article created successfully", zap.String("id", article.ID), zap.String("tenantId", *tenantID))
	return article, nil
}

// GetAllArticles retrieves all articles for a tenant with pagination
func (s *ArticleService) GetAllArticles(tenantID string, page, pageSize int, articleType *string) (*dto.ArticlePaginationResponse, error) {
	articles, total, err := s.Repo.FindAllWithPagination(tenantID, page, pageSize, articleType)
	if err != nil {
		logger.Log.Error("Failed to fetch all articles with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, err
	}

	articleResponses := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponses[i] = dto.ArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			Slug:        article.Slug,
			CoverImage:  article.CoverImage,
			Body:        article.Body,
			Excerpt:     article.Excerpt,
			IsPublished: article.IsPublished,
			ArticleType: article.ArticleType,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 && total > 0 { // Untuk memastikan 1 halaman jika ada data tapi pageSize sangat besar
		totalPages = 1
	}

	response := &dto.ArticlePaginationResponse{
		Data:       articleResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	logger.Log.Info("Fetched articles with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Int64("total", total))
	return response, nil
}

// GetArticleByID retrieves an article by its ID
func (s *ArticleService) GetArticleByID(id string, tenantID string) (*models.Article, error) {
	article, err := s.Repo.FindByID(id, tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("Article not found", zap.String("id", id), zap.String("tenantId", tenantID))
			return nil, errors.New("artikel tidak ditemukan")
		}
		logger.Log.Error("Failed to get article by ID", zap.String("id", id), zap.String("tenantId", tenantID), zap.Error(err))
		return nil, err
	}
	return article, nil
}

// UpdateArticle updates an existing article
func (s *ArticleService) UpdateArticle(id string, updateDTO dto.UpdateArticleDTO, tenantID string) (*models.Article, error) {
	article, err := s.Repo.FindByID(id, tenantID)
	if err != nil {
		logger.Log.Warn("Failed to find article for update", zap.String("id", id), zap.String("tenantId", tenantID), zap.Error(err))
		return nil, err
	}

	// Update fields if provided
	if updateDTO.Title != nil {
		article.Title = *updateDTO.Title
	}
	if updateDTO.Slug != nil {
		// Check for duplicate slug if changed
		if article.Slug != *updateDTO.Slug {
			existingArticle, err := s.Repo.FindBySlug(*updateDTO.Slug, tenantID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Log.Error("Failed to check existing slug during update", zap.String("slug", *updateDTO.Slug), zap.String("tenantId", tenantID), zap.Error(err))
				return nil, err
			}
			if existingArticle != nil && existingArticle.ID != article.ID {
				return nil, errors.New("slug sudah ada untuk tenant ini")
			}
		}
		article.Slug = *updateDTO.Slug
	}
	if updateDTO.CoverImage != nil {
		article.CoverImage = updateDTO.CoverImage
	}
	if updateDTO.Body != nil {
		article.Body = updateDTO.Body
	}
	if updateDTO.Excerpt != nil {
		article.Excerpt = updateDTO.Excerpt
	}
	if updateDTO.IsPublished != nil {
		article.IsPublished = *updateDTO.IsPublished
	}
	if updateDTO.ArticleType != nil {
		article.ArticleType = updateDTO.ArticleType
	}

	if err := s.Repo.Update(article); err != nil {
		logger.Log.Error("Failed to update article", zap.String("id", id), zap.String("tenantId", tenantID), zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Article updated successfully", zap.String("id", id), zap.String("tenantId", tenantID))
	return article, nil
}

// DeleteArticle deletes an article by its ID
func (s *ArticleService) DeleteArticle(id string, tenantID string) error {
	if err := s.Repo.Delete(id, tenantID); err != nil {
		logger.Log.Error("Failed to delete article", zap.String("id", id), zap.String("tenantId", tenantID), zap.Error(err))
		return err
	}
	logger.Log.Info("Article deleted successfully", zap.String("id", id), zap.String("tenantId", tenantID))
	return nil
}

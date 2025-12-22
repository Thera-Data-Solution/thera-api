// repositories/article_repository.go
package repositories

import (
	"thera-api/models" // Sesuaikan path jika perlu

	"gorm.io/gorm"
)

type ArticleRepository struct {
	DB *gorm.DB
}

// NewArticleRepository creates a new instance of ArticleRepository
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{DB: db}
}

// Create creates a new Article record in the database
func (r *ArticleRepository) Create(article *models.Article) error {
	return r.DB.Create(article).Error
}

// FindAllWithPagination retrieves a list of articles for a given tenant with pagination
func (r *ArticleRepository) FindAllWithPagination(tenantID string, page, pageSize int, articleType *string) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.DB.Model(&models.Article{}).Where("tenant_id = ?", tenantID)

	if articleType != nil && *articleType != "" {
		query = query.Where("article_type = ?", *articleType)
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan pagination
	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC"). // Urutkan berdasarkan tanggal terbaru
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindByID retrieves an Article by its ID and tenantID
func (r *ArticleRepository) FindByID(id string, tenantID string) (*models.Article, error) {
	var article models.Article
	err := r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindBySlug retrieves an Article by its Slug and tenantID
func (r *ArticleRepository) FindBySlug(slug string, tenantID string) (*models.Article, error) {
	var article models.Article
	err := r.DB.Where("slug = ? AND tenant_id = ?", slug, tenantID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Update updates an existing Article record in the database
func (r *ArticleRepository) Update(article *models.Article) error {
	return r.DB.Save(article).Error
}

// Delete deletes an Article record from the database by its ID and tenantID
func (r *ArticleRepository) Delete(id string, tenantID string) error {
	return r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Article{}).Error
}

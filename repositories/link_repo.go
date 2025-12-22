package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type LinkRepository struct {
	DB *gorm.DB
}

func (r *LinkRepository) Create(link *models.Link) error {
	return r.DB.Create(link).Error
}

func (r *LinkRepository) FindAll(tenant string) ([]models.Link, error) {
	var links []models.Link
	err := r.DB.Where("tenant_id = ?", tenant).Order("\"order\" ASC").Find(&links).Error
	return links, err
}

func (r *LinkRepository) FindByID(id string, tenant string) (*models.Link, error) {
	var link models.Link
	err := r.DB.First(&link, "id = ? AND tenant_id = ?", id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Update(link *models.Link) error {
	return r.DB.Save(link).Error
}

func (r *LinkRepository) Delete(id string, tenant string) error {
	return r.DB.Delete(&models.Link{}, "id = ? AND tenant_id = ?", id, tenant).Error
}

func (r *LinkRepository) GetMaxOrderByTenantID(tenantID string) (int, error) {
	var maxOrder int
	err := r.DB.Model(&models.Link{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Row().
		Scan(&maxOrder)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return maxOrder, nil
}

func (r *LinkRepository) GetLinkByOrder(order int, tenantID string) (*models.Link, error) {
	var link models.Link
	err := r.DB.Where("order = ? AND tenant_id = ?", order, tenantID).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) ReorderLinksBatch(tenantID string) error {
	var links []models.Link
	// Ambil semua link, urutkan berdasarkan order saat ini
	if err := r.DB.Where("tenant_id = ?", tenantID).Order("\"order\" ASC").Find(&links).Error; err != nil {
		return err
	}

	tx := r.DB.Begin() // Memulai transaksi
	if tx.Error != nil {
		return tx.Error
	}

	for i, link := range links {
		newOrder := i + 1 // Order baru akan 1, 2, 3, ...
		// Perbarui hanya jika order berubah untuk menghindari operasi DB yang tidak perlu
		if link.Order == nil || *link.Order != newOrder {
			if err := tx.Model(&models.Link{}).Where("id = ?", link.ID).Update("order", newOrder).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error // Commit transaksi
}

func (r *LinkRepository) FindLinksBetweenOrders(tenantID string, startOrder, endOrder int) ([]models.Link, error) {
	var links []models.Link
	err := r.DB.
		Where("tenant_id = ? AND \"order\" BETWEEN ? AND ?", tenantID, startOrder, endOrder).
		Order("\"order\" ASC").
		Find(&links).Error
	return links, err
}

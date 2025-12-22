package services

import (
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LinkService struct {
	Repo *repositories.LinkRepository
}

func (s *LinkService) GetAllLinks(tenantId string) ([]models.Link, error) {
	links, err := s.Repo.FindAll(tenantId)
	if err != nil {
		logger.Log.Error("Failed to fetch all links", zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Fetched all links", zap.String("tenantId", tenantId), zap.Int("count", len(links)))
	return links, nil
}

func (s *LinkService) GetLinkByID(id string, tenantId string) (*models.Link, error) {
	link, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		logger.Log.Warn("Link not found", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	return link, nil
}

func (s *LinkService) CreateLink(name, value, linkType string, icon *string, tenantId *string) (*models.Link, error) {
	maxOrder, err := s.Repo.GetMaxOrderByTenantID(*tenantId)
	if err != nil {
		logger.Log.Error("Failed to get max order for tenant", zap.String("tenantId", *tenantId), zap.Error(err))
		return nil, err
	}
	newOrder := maxOrder + 1
	link := &models.Link{
		Name:     name,
		Value:    value,
		Type:     linkType,
		Icon:     icon,
		Order:    &newOrder,
		TenantId: tenantId,
	}
	if err := s.Repo.Create(link); err != nil {
		logger.Log.Error("Failed to create link", zap.String("name", name), zap.String("type", linkType), zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Link created successfully", zap.String("id", link.ID), zap.String("tenantId", *tenantId))
	return link, nil
}

func (s *LinkService) UpdateLink(id string, dto map[string]interface{}, tenantId string) (*models.Link, error) {
	link, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		logger.Log.Warn("Failed to find link for update", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}

	if name, ok := dto["name"].(string); ok {
		link.Name = name
	}
	if value, ok := dto["value"].(string); ok {
		link.Value = value
	}
	if linkType, ok := dto["type"].(string); ok {
		link.Type = linkType
	}
	if icon, ok := dto["icon"].(string); ok {
		link.Icon = &icon
	}
	if order, ok := dto["order"].(int); ok {
		link.Order = &order
	}

	if err := s.Repo.Update(link); err != nil {
		logger.Log.Error("Failed to update link", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Link updated successfully", zap.String("id", id), zap.String("tenantId", tenantId))
	return link, nil
}

func (s *LinkService) DeleteLink(id string, tenantId string) error {
	if err := s.Repo.Delete(id, tenantId); err != nil {
		logger.Log.Error("Failed to delete link", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return err
	}
	logger.Log.Info("Link deleted successfully", zap.String("id", id), zap.String("tenantId", tenantId))
	return nil
}

func (s *LinkService) UpdateOrder(id string, newOrder int, tenantId string) (*models.Link, error) {
	link, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}

	// Ambil semua link tenant
	links, err := s.Repo.FindAll(tenantId)
	if err != nil {
		return nil, err
	}

	// Perbarui order link lain supaya tidak duplikat
	for _, l := range links {
		if l.ID == link.ID {
			continue
		}
		if l.Order != nil && *l.Order >= newOrder {
			temp := *l.Order + 1
			l.Order = &temp
			_ = s.Repo.Update(&l)
		}
	}

	link.Order = &newOrder
	if err := s.Repo.Update(link); err != nil {
		return nil, err
	}

	return link, nil
}
func (s *LinkService) reorderLink(linkToMove *models.Link, newOrder int, tenantID string) error {
	currentOrder := *linkToMove.Order

	if newOrder == currentOrder {
		return nil // Tidak ada perubahan posisi
	}

	// Ambil order max saat ini untuk tenant ini
	maxOrder, err := s.Repo.GetMaxOrderByTenantID(tenantID)
	if err != nil {
		return err
	}

	// Validasi newOrder
	if newOrder < 1 || newOrder > maxOrder {
		if newOrder == maxOrder+1 && currentOrder == maxOrder {
			// Kasus khusus: jika memindahkan link paling bawah ke satu posisi lebih bawah
			// Ini biasanya tidak terjadi dengan MoveDown, karena maxOrder akan sama
			// Tapi untuk UpdateLink, ini bisa terjadi jika maxOrder + 1 adalah target.
			// Biarkan ini sebagai catatan, logika dibawah sudah cukup.
		} else if newOrder > maxOrder && currentOrder < maxOrder {
			// Jika user mencoba mengatur order terlalu tinggi dari maxOrder saat ini
			newOrder = maxOrder // Sesuaikan menjadi maxOrder jika melebihi
		} else if newOrder < 1 {
			newOrder = 1 // Sesuaikan menjadi 1 jika kurang dari 1
		}
	}

	// Mulai transaksi untuk memastikan konsistensi
	tx := s.Repo.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if currentOrder < newOrder {
		// Bergerak ke bawah: geser semua link di antara (currentOrder + 1) dan newOrder ke atas (kurangi 1)
		// Misal: A(1) B(2) C(3) D(4) E(5). Pindahkan A ke 3.
		// A(1) -> A(temp)
		// B(2) -> B(1)
		// C(3) -> C(2)
		// D(4) -> D(3)
		// E(5) tetap
		// A(temp) -> A(3)
		if err := tx.Model(&models.Link{}).
			Where("tenant_id = ? AND \"order\" > ? AND \"order\" <= ?", tenantID, currentOrder, newOrder).
			Update("order", gorm.Expr("\"order\" - ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else { // currentOrder > newOrder
		// Bergerak ke atas: geser semua link di antara newOrder dan (currentOrder - 1) ke bawah (tambah 1)
		// Misal: A(1) B(2) C(3) D(4) E(5). Pindahkan E ke 2.
		// E(5) -> E(temp)
		// A(1) tetap
		// B(2) -> B(3)
		// C(3) -> C(4)
		// D(4) -> D(5)
		// E(temp) -> E(2)
		if err := tx.Model(&models.Link{}).
			Where("tenant_id = ? AND \"order\" >= ? AND \"order\" < ?", tenantID, newOrder, currentOrder).
			Update("order", gorm.Expr("\"order\" + ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update link yang dipindahkan ke posisi barunya
	if err := tx.Model(linkToMove).Update("order", newOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *LinkService) MoveUp(id string, tenantId string) (*models.Link, error) {
	link, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}
	if link.Order == nil || *link.Order <= 1 {
		return link, nil // sudah paling atas atau order belum diinisialisasi
	}

	newOrder := *link.Order - 1
	err = s.reorderLink(link, newOrder, tenantId)
	if err != nil {
		logger.Log.Error("Failed to move link up", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}

	// Refresh link dari DB untuk mendapatkan order terbaru setelah reorder
	updatedLink, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		logger.Log.Error("Failed to fetch updated link after move up", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	return updatedLink, nil
}

func (s *LinkService) MoveDown(id string, tenantId string) (*models.Link, error) {
	link, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}

	maxOrder, err := s.Repo.GetMaxOrderByTenantID(tenantId)
	if err != nil {
		logger.Log.Error("Failed to get max order for tenant during MoveDown", zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}

	if link.Order == nil || *link.Order >= maxOrder {
		return link, nil // sudah paling bawah atau order belum diinisialisasi
	}

	newOrder := *link.Order + 1
	err = s.reorderLink(link, newOrder, tenantId)
	if err != nil {
		logger.Log.Error("Failed to move link down", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}

	// Refresh link dari DB untuk mendapatkan order terbaru setelah reorder
	updatedLink, err := s.Repo.FindByID(id, tenantId)
	if err != nil {
		logger.Log.Error("Failed to fetch updated link after move down", zap.String("id", id), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	return updatedLink, nil
}

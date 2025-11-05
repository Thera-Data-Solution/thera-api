package repositories

import (
	"fmt"
	"thera-api/models"

	"gorm.io/gorm"
)

type GalleryRepository struct {
	DB *gorm.DB
}

func (r *GalleryRepository) FindAll(tenant string) ([]models.Gallery, error) {
	fmt.Println("testing disini3")
	var gallery []models.Gallery
	err := r.DB.Where(`tenant_id = ?`, tenant).Find(&gallery).Error
	return gallery, err
}

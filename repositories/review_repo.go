package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func (r *ReviewRepository) Upsert(review *models.Review) error {
	return r.DB.Save(review).Error
}

func (r *ReviewRepository) FindByBookingID(bookingId string) (*models.Review, error) {
	var review models.Review
	err := r.DB.Where("booking_id = ?", bookingId).First(&review).Error
	return &review, err
}

func (r *ReviewRepository) GetPublicReviews(tenantId string) ([]models.Booked, error) {
	var results []models.Booked
	err := r.DB.
		Preload("User").
		Preload("Schedule.Categories").
		Preload("Event").
		Preload("Review").
		Joins("JOIN reviews ON reviews.booking_id = bookeds.id").
		Where("bookeds.tenant_id = ? AND reviews.is_approved = ?", tenantId, true).
		Find(&results).Error
	return results, err
}

func (r *ReviewRepository) GetUserHistory(userId, tenantId string) ([]models.Booked, error) {
	var results []models.Booked
	err := r.DB.
		Preload("Schedule.Categories").
		Preload("Event").
		Preload("Review").
		Where("user_id = ? AND tenant_id = ?", userId, tenantId).
		Order("booked_at DESC").
		Find(&results).Error
	return results, err
}

func (r *ReviewRepository) GetAllForAdmin(tenantId string) ([]models.Booked, error) {
	var results []models.Booked
	err := r.DB.
		Preload("User").
		Preload("Schedule.Categories").
		Preload("Event").
		Preload("Review").
		Joins("JOIN reviews ON reviews.booking_id = bookeds.id").
		Where("bookeds.tenant_id = ?", tenantId).
		Find(&results).Error
	return results, err
}

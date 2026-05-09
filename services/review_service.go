package services

import (
	"errors"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"

	"go.uber.org/zap"
)

type ReviewService struct {
	ReviewRepo  *repositories.ReviewRepository
	BookingRepo *repositories.BookedRepository
}

func (s *ReviewService) SubmitReview(userId string, req dto.ReviewUpsertRequest, tenantId string) error {
	booking, err := s.BookingRepo.GetById(req.BookingId, tenantId)
	if err != nil || booking.UserId != userId {
		return errors.New("booking tidak ditemukan atau bukan milik Anda")
	}

	var targetId, targetType string
	if booking.ScheduleId != nil {
		targetId = booking.Schedule.CategoryId
		targetType = "category"
	} else if booking.EventId != nil {
		targetId = *booking.EventId
		targetType = "event"
	}

	review := &models.Review{
		BookingId:   req.BookingId,
		UserId:      userId,
		TargetId:    targetId,
		TargetType:  targetType,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
		IsApproved:  false,
	}

	existing, _ := s.ReviewRepo.FindByBookingID(req.BookingId)
	if existing != nil {
		review.ID = existing.ID
	}

	return s.ReviewRepo.Upsert(review)
}
func (s *ReviewService) GetHistory(userId, tenantId string) ([]dto.UserHistoryResponse, error) {
	data, err := s.ReviewRepo.GetUserHistory(userId, tenantId)
	if err != nil {
		return nil, err
	}

	var res []dto.UserHistoryResponse
	for _, b := range data {
		item := dto.UserHistoryResponse{
			BookingId: b.ID,
			Status:    b.Status,
			BookedAt:  b.BookedAt,
		}

		// Cek Schedule aman
		if b.ScheduleId != nil && b.Schedule != nil {
			item.ItemName = b.Schedule.Categories.Name
			item.ScheduleDate = b.Schedule.DateTime
			if b.Schedule.Categories.Image != nil {
				item.ItemImage = *b.Schedule.Categories.Image
			}
		} else if b.EventId != nil && b.Event != nil {
			item.ItemName = b.Event.Name
			item.ScheduleDate = b.Event.StartAt
			if b.Event.Image != nil {
				item.ItemImage = *b.Event.Image
			}
		}

		// Cek Review aman
		if b.Review != nil {
			item.Review = &dto.ReviewInfo{
				Content:     b.Review.Content,
				IsAnonymous: b.Review.IsAnonymous,
				// IsApproved tidak perlu di-set jika tidak ada di DTO,
				// Go akan otomatis menggunakan default value (false) jika field ada di struct.
			}
		}

		res = append(res, item)
	}
	return res, nil
}

func (s *ReviewService) GetPublicTestimonials(tenantId string) ([]dto.PublicReviewResponse, error) {
	logger.Log.Info("Fetching public testimonials", zap.String("tenantId", tenantId))
	data, err := s.ReviewRepo.GetPublicReviews(tenantId)
	if err != nil {
		return nil, err
	}

	res := make([]dto.PublicReviewResponse, 0)
	for _, b := range data {
		userName := b.User.FullName
		if b.Review.IsAnonymous {
			runes := []rune(userName)
			if len(runes) > 0 {
				userName = string(runes[0]) + "*****"
			}
		}

		targetName := ""
		if b.ScheduleId != nil {
			targetName = b.Schedule.Categories.Name
		} else {
			targetName = b.Event.Name
		}

		res = append(res, dto.PublicReviewResponse{
			UserName: userName,
			Content:  b.Review.Content,
			Type:     b.Review.TargetType,
			Target:   targetName,
			Date:     b.Review.CreatedAt,
		})
	}
	return res, nil
}

func (s *ReviewService) AdminGetAllReviews(tenantId string) ([]dto.AdminReviewResponse, error) {
	// Kita gunakan fungsi GetAllForAdmin yang melakukan preload User, Schedule, dan Event
	data, err := s.ReviewRepo.GetAllForAdmin(tenantId)
	if err != nil {
		return nil, err
	}

	res := make([]dto.AdminReviewResponse, 0)
	for _, b := range data {
		if b.Review == nil {
			continue
		}

		targetName := ""
		if b.ScheduleId != nil {
			targetName = b.Schedule.Categories.Name
		} else if b.EventId != nil {
			targetName = b.Event.Name
		}

		res = append(res, dto.AdminReviewResponse{
			ID:          b.Review.ID,
			BookingId:   b.ID,
			UserName:    b.User.FullName,
			TargetName:  targetName,
			TargetType:  b.Review.TargetType,
			Content:     b.Review.Content,
			IsApproved:  b.Review.IsApproved,
			IsAnonymous: b.Review.IsAnonymous,
			CreatedAt:   b.Review.CreatedAt,
		})
	}
	return res, nil
}

func (s *ReviewService) AdminUpdateReview(reviewId string, req dto.AdminUpdateReviewRequest) error {
	review, err := s.ReviewRepo.FindByID(reviewId)
	if err != nil {
		return errors.New("ulasan tidak ditemukan")
	}

	// Update content jika dikirim
	if req.Content != nil {
		review.Content = *req.Content
	}

	// Update status approval jika dikirim
	if req.IsApproved != nil {
		review.IsApproved = *req.IsApproved
	}

	return s.ReviewRepo.Update(review)
}

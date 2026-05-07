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

		if b.ScheduleId != nil {
			item.ItemName = b.Schedule.Categories.Name
			item.ItemImage = *b.Schedule.Categories.Image
			item.ScheduleDate = b.Schedule.DateTime
		} else if b.EventId != nil {
			item.ItemName = b.Event.Name
			item.ItemImage = *b.Event.Image
			item.ScheduleDate = b.Event.StartAt
		}

		if b.Review != nil {
			item.Review = &dto.ReviewInfo{
				Content:    b.Review.Content,
				IsApproved: b.Review.IsApproved,
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
			userName = "Anonymous"
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

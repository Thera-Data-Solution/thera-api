package services

import (
	"errors"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"go.uber.org/zap"
)

type BookedService struct {
	BookingRepo  *repositories.BookedRepository
	ScheduleRepo *repositories.SchedulesRepository
}

func NewBookedService(bookingRepo *repositories.BookedRepository, scheduleRepo *repositories.SchedulesRepository) *BookedService {
	return &BookedService{
		BookingRepo:  bookingRepo,
		ScheduleRepo: scheduleRepo,
	}
}

func (s *BookedService) Create(userId, scheduleId string, tenantId string) error {
	logger.Log.Info("Create booking called", zap.String("userId", userId), zap.String("scheduleId", scheduleId), zap.String("tenantId", tenantId))

	schedule, err := s.ScheduleRepo.FindByID(scheduleId, tenantId)
	if err != nil {
		logger.Log.Warn("Schedule tidak ditemukan", zap.String("scheduleId", scheduleId))
		return err
	}

	if schedule.Status != "ENABLE" {
		logger.Log.Warn("Schedule tidak tersedia untuk dibooking", zap.String("scheduleId", scheduleId))
		return errors.New("jadwal tidak tersedia untuk dibooking")
	}

	schedule.Status = "BOOKED"
	if err := s.ScheduleRepo.Update(schedule); err != nil {
		logger.Log.Error("Gagal update status schedule", zap.String("scheduleId", scheduleId), zap.Error(err))
		return err
	}

	booked := &models.Booked{
		UserId:     userId,
		ScheduleId: scheduleId,
		BookedAt:   time.Now(),
		TenantId:   &tenantId,
	}

	if err := s.BookingRepo.Create(booked); err != nil {
		logger.Log.Error("Gagal membuat booking", zap.String("userId", userId), zap.String("scheduleId", scheduleId), zap.Error(err))
		return err
	}

	logger.Log.Info("Booking berhasil dibuat", zap.String("userId", userId), zap.String("scheduleId", scheduleId))
	return nil
}

func (s *BookedService) GetAll(tenantId string) ([]models.Booked, error) {
	bookings, err := s.BookingRepo.GetAll(tenantId)
	if err != nil {
		logger.Log.Error("Gagal fetch booking", zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	return bookings, nil
}

func (s *BookedService) GetByUser(tenantId string, userId string) ([]models.Booked, error) {
	bookings, err := s.BookingRepo.GetByUser(tenantId, userId)
	if err != nil {
		logger.Log.Error("Gagal fetch booking user", zap.String("userId", userId), zap.String("tenantId", tenantId), zap.Error(err))
		return nil, err
	}
	return bookings, nil
}

func (s *BookedService) GetById(id string, tenantId string) (*models.Booked, error) {
	booked, err := s.BookingRepo.GetById(id, tenantId)
	if err != nil {
		logger.Log.Warn("Booking tidak ditemukan", zap.String("id", id), zap.String("tenantId", tenantId))
		return nil, err
	}
	return booked, nil
}

func (s *BookedService) Update(booked *models.Booked) error {
	if err := s.BookingRepo.Update(booked); err != nil {
		logger.Log.Error("Gagal update booking", zap.String("bookingId", booked.ID), zap.Error(err))
		return err
	}
	logger.Log.Info("Booking berhasil diupdate", zap.String("bookingId", booked.ID))
	return nil
}

func (s *BookedService) Delete(id string, tenantId string) error {
	if err := s.BookingRepo.Delete(id, tenantId); err != nil {
		logger.Log.Error("Gagal delete booking", zap.String("bookingId", id), zap.Error(err))
		return err
	}
	logger.Log.Info("Booking berhasil dihapus", zap.String("bookingId", id))
	return nil
}

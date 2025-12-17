package services

import (
	"errors"
	"fmt"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
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

func (s *BookedService) Create(userId, scheduleId string, tenantId string, customAnswer datatypes.JSON) error {
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
		UserId:       userId,
		ScheduleId:   scheduleId,
		BookedAt:     time.Now(),
		TenantId:     &tenantId,
		CustomAnswer: customAnswer,
	}

	if err := s.BookingRepo.Create(booked); err != nil {
		logger.Log.Error("Gagal membuat booking", zap.String("userId", userId), zap.String("scheduleId", scheduleId), zap.Error(err))
		return err
	}

	logger.Log.Info("Booking berhasil dibuat", zap.String("userId", userId), zap.String("scheduleId", scheduleId))
	return nil
}

func (s *BookedService) GetAll(tenantId string, limit, offset int) ([]models.Booked, int64, error) {
	bookings, total, err := s.BookingRepo.GetAll(tenantId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
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

func (s *BookedService) AddTestimoni(id string, testimoni *string, tenantId string) (*models.Booked, error) {
	logger.Log.Info("Add testimoni called", zap.String("id", id), zap.String("tenantId", tenantId))

	booked, err := s.BookingRepo.GetById(id, tenantId)
	if err != nil {
		return nil, err
	}

	if booked.Schedule.Status != "CLOSED" {
		logger.Log.Warn("Status jadwal tidak valid", zap.String("scheduleId", booked.ScheduleId))
		fmt.Println(booked.Schedule.Status)
		return nil, errors.New("testimoni hanya dapat diisi jika status jadwal sudah selesai")
	}

	booked.Testimoni = testimoni
	if err := s.BookingRepo.Update(booked); err != nil {
		return nil, err
	}

	return booked, nil
}

func (s *BookedService) Cancel(bookingId, tenantId string) error {
	logger.Log.Info("Cancel booking called", zap.String("bookingId", bookingId), zap.String("tenantId", tenantId))

	booking, err := s.BookingRepo.GetById(bookingId, tenantId)
	if err != nil {
		logger.Log.Warn("Booking tidak ditemukan", zap.String("bookingId", bookingId))
		return err
	}
	schedule, err := s.ScheduleRepo.FindByID(booking.ScheduleId, tenantId)
	if err != nil {
		logger.Log.Warn("Schedule terkait tidak ditemukan", zap.String("scheduleId", booking.ScheduleId))
		return err
	}

	schedule.Status = "ENABLE"
	if err := s.ScheduleRepo.Update(schedule); err != nil {
		logger.Log.Error("Gagal mengubah status schedule", zap.String("scheduleId", booking.ScheduleId), zap.Error(err))
		return err
	}

	if err := s.BookingRepo.Delete(bookingId, tenantId); err != nil {
		logger.Log.Error("Gagal menghapus booking", zap.String("bookingId", bookingId), zap.Error(err))
		return err
	}

	logger.Log.Info("Booking berhasil dibatalkan", zap.String("bookingId", bookingId), zap.String("scheduleId", booking.ScheduleId))
	return nil
}

func (s *BookedService) GetAllTestimoni(tenantId string) ([]models.Booked, error) {
	logger.Log.Info("Get all testimoni called", zap.String("tenantId", tenantId))

	return s.BookingRepo.GetAllWithTestimoni(tenantId)
}

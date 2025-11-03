package services

import (
	"errors"
	"time"

	"thera-api/models"
	"thera-api/repositories"
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
	// 1️⃣ Ambil data schedule berdasarkan ID + tenant
	schedule, err := s.ScheduleRepo.FindByID(scheduleId, tenantId)
	if err != nil {
		return err
	}

	// 2️⃣ Cek status schedule
	if schedule.Status != "ENABLE" {
		return errors.New("jadwal tidak tersedia untuk dibooking")
	}

	// 3️⃣ Ubah status schedule jadi BOOKED
	schedule.Status = "BOOKED"
	if err := s.ScheduleRepo.Update(schedule); err != nil {
		return err
	}

	// 4️⃣ Buat booking baru
	booked := &models.Booked{
		UserId:     userId,
		ScheduleId: scheduleId,
		BookedAt:   time.Now(),
		TenantId:   &tenantId,
	}

	return s.BookingRepo.Create(booked)
}

func (s *BookedService) GetAll(tenantId string) ([]models.Booked, error) {
	return s.BookingRepo.GetAll(tenantId)
}

func (s *BookedService) GetByUser(tenantId string, userId string) ([]models.Booked, error) {
	return s.BookingRepo.GetByUser(tenantId, userId)
}

func (s *BookedService) GetById(id string, tenantId string) (*models.Booked, error) {
	return s.BookingRepo.GetById(id, tenantId)
}

func (s *BookedService) Update(booked *models.Booked) error {
	return s.BookingRepo.Update(booked)
}

func (s *BookedService) Delete(id string, tenantId string) error {
	return s.BookingRepo.Delete(id, tenantId)
}

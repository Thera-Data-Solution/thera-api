package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"thera-api/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type BookedService struct {
	BookingRepo  *repositories.BookedRepository
	ScheduleRepo *repositories.SchedulesRepository
	SettingRepo  repositories.SettingRepo
	EventRepo    *repositories.EventsRepository
}

func NewBookedService(bookingRepo *repositories.BookedRepository, scheduleRepo *repositories.SchedulesRepository, settingRepo repositories.SettingRepo, eventRepo *repositories.EventsRepository) *BookedService {
	return &BookedService{
		BookingRepo:  bookingRepo,
		ScheduleRepo: scheduleRepo,
		SettingRepo:  settingRepo,
		EventRepo:    eventRepo,
	}
}

func (s *BookedService) Create(userId, targetId string, tenantId string, customAnswer datatypes.JSON, bookingType int) error {
	logger.Log.Info("Create booking called", zap.String("userId", userId), zap.Int("type", bookingType))

	var booked models.Booked
	booked.UserId = userId
	booked.TenantId = &tenantId
	booked.CustomAnswer = customAnswer
	booked.Status = "PENDING"

	switch bookingType {
	case 1:
		if existing, err := s.BookingRepo.GetLatestByUserAndSchedule(userId, targetId, tenantId); err == nil && existing != nil {
			return errors.New("Anda sudah melakukan booking untuk jadwal ini")
		}
		schedule, err := s.ScheduleRepo.FindByID(targetId, tenantId)
		if err != nil {
			return errors.New("jadwal tidak ditemukan")
		}
		if schedule.Status != "ENABLE" && schedule.Status != "AVAILABLE" {
			return errors.New("jadwal tidak tersedia")
		}

		if !schedule.Categories.IsGroup {
			schedule.Status = "BOOKED"
			s.ScheduleRepo.Update(schedule)
		}

		booked.ScheduleId = &targetId
	case 2:
		if existing, err := s.BookingRepo.GetLatestByUserAndEvent(userId, targetId, tenantId); err == nil && existing != nil {
			return errors.New("Anda sudah melakukan booking untuk event ini")
		}
		event, err := s.EventRepo.FindByID(targetId, tenantId) // Pastikan EventRepo sudah ada
		if err != nil {
			return errors.New("event tidak ditemukan")
		}
		if event.Status != "available" {
			return errors.New("event sudah tidak tersedia")
		}

		booked.EventId = &targetId
	default:
		return errors.New("tipe booking tidak valid")
	}
	if err := s.BookingRepo.Create(&booked); err != nil {
		return err
	}

	bookedFull, err := s.BookingRepo.GetByIDWithDetails(booked.ID, tenantId)
	if err == nil {
		go s.sendTelegramNotificationAsync(bookedFull, tenantId)
	}
	return nil
}

func (s *BookedService) AdminChangeStatus(bookingId string, newStatus string, tenantId string) error {
	logger.Log.Info("Admin change booking status", zap.String("id", bookingId), zap.String("status", newStatus))

	booking, err := s.BookingRepo.GetByIDWithDetails(bookingId, tenantId)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	if newStatus == "CANCELLED" {
		if booking.ScheduleId != nil {
			schedule, err := s.ScheduleRepo.FindByID(*booking.ScheduleId, tenantId)
			if err == nil && !schedule.Categories.IsGroup {
				schedule.Status = "ENABLE"
				s.ScheduleRepo.Update(schedule)
			}
		}
	}

	if newStatus == "CONFIRMED" {
		if booking.ScheduleId != nil {
			schedule, err := s.ScheduleRepo.FindByID(*booking.ScheduleId, tenantId)
			if err == nil && !schedule.Categories.IsGroup {
				schedule.Status = "BOOKED"
				s.ScheduleRepo.Update(schedule)
			}
		}
	}

	booking.Status = newStatus
	return s.BookingRepo.Update(booking)
}

func (s *BookedService) CloseSchedule(scheduleId string, tenantId string) error {
	schedule, err := s.ScheduleRepo.FindByID(scheduleId, tenantId)
	if err != nil {
		return errors.New("jadwal tidak ditemukan")
	}

	schedule.Status = "CLOSED"
	if err := s.ScheduleRepo.Update(schedule); err != nil {
		return err
	}
	return s.BookingRepo.UpdateStatusBySchedule(scheduleId, "CONFIRMED", "CLOSED")
}

func (s *BookedService) GetAllForAdmin(tenantId string, limit, offset int) ([]dto.GetAllBookingResponse, int64, error) {
	data, total, err := s.BookingRepo.GetAll(tenantId, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.GetAllBookingResponse
	for _, b := range data {
		var customAnswer []dto.CustomAnswerResponse
		if err := json.Unmarshal(b.CustomAnswer, &customAnswer); err != nil {
			logger.Log.Error("Gagal unmarshal customAnswer", zap.String("bookingId", b.ID), zap.Error(err))
		}
		res := dto.GetAllBookingResponse{
			ID:           b.ID,
			Avatar:       b.User.Avatar,
			UserId:       b.UserId,
			FullName:     b.User.FullName,
			Email:        b.User.Email,
			Phone:        b.User.Phone,
			BookAt:       b.BookedAt,
			Status:       b.Status,
			CustomAnswer: customAnswer,
		}
		if b.ScheduleId != nil && b.Schedule != nil {
			res.ScheduleId = b.ScheduleId
			name := b.Schedule.Categories.Name
			res.ScheduleName = &name
			res.ScheduleImage = b.Schedule.Categories.Image
			res.ScheduleDate = &b.Schedule.DateTime
		}

		if b.EventId != nil && b.Event != nil {
			res.EventId = b.EventId
			name := b.Event.Name
			res.EventName = &name
			res.EventImage = b.Event.Image
			res.EventDate = &b.Event.StartAt
		}

		response = append(response, res)
	}

	return response, total, nil
}

func (s *BookedService) Update(booked *models.Booked) error {
	if err := s.BookingRepo.Update(booked); err != nil {
		logger.Log.Error("Gagal update booking", zap.String("bookingId", booked.ID), zap.Error(err))
		return err
	}
	logger.Log.Info("Booking berhasil diupdate", zap.String("bookingId", booked.ID))
	return nil
}

func (s *BookedService) Cancel(
	targetId string,
	bookingType int,
	userId *string,
	tenantId string,
) error {
	if userId == nil {
		return errors.New("user ID tidak valid")
	}

	logger.Log.Info("Cancel booking request", zap.String("targetId", targetId), zap.Int("type", bookingType), zap.String("userId", *userId))

	var booking *models.Booked
	var err error

	switch bookingType {
	case 1:
		booking, err = s.BookingRepo.GetLatestByUserAndScheduleWithoutStatus(*userId, targetId, tenantId)
		if err != nil {
			return errors.New("data booking jadwal tidak ditemukan")
		}
		if booking.Status == "CANCELLED" || booking.Status == "CLOSED" {
			return fmt.Errorf("booking tidak dapat dibatalkan karena status sudah %s", booking.Status)
		}
		schedule, err := s.ScheduleRepo.FindByID(*booking.ScheduleId, tenantId)
		if err == nil && schedule != nil {
			if !schedule.Categories.IsGroup {
				schedule.Status = "ENABLE"
				if err := s.ScheduleRepo.Update(schedule); err != nil {
					logger.Log.Error("Gagal mengupdate status schedule ke ENABLE", zap.Error(err))
				}
			}
		}

	case 2:
		booking, err = s.BookingRepo.GetLatestByUserAndEventWithoutStatus(*userId, targetId, tenantId)
		if err != nil {
			return errors.New("data booking event tidak ditemukan")
		}

		if booking.Status == "CANCELLED" || booking.Status == "CLOSED" {
			return fmt.Errorf("booking tidak dapat dibatalkan karena status sudah %s", booking.Status)
		}
	default:
		return errors.New("tipe booking tidak valid")
	}
	booking.Status = "CANCELLED"
	if err := s.BookingRepo.Update(booking); err != nil {
		logger.Log.Error("Gagal memperbarui status booking menjadi CANCELLED", zap.String("bookingId", booking.ID), zap.Error(err))
		return err
	}

	logger.Log.Info("Booking berhasil dibatalkan", zap.String("bookingId", booking.ID), zap.String("targetId", targetId))
	return nil
}

// Contoh untuk Telegram (berlaku sama untuk Discord)
func (s *BookedService) sendTelegramNotificationAsync(booked *models.Booked, tenantId string) {
	logger.Log.Info("Telegram async started", zap.String("bookingId", booked.ID))

	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered from panic in Telegram notification", zap.Any("error", r))
		}
	}()

	setting, err := s.SettingRepo.FindByTenantId(tenantId)
	if err != nil {
		logger.Log.Warn("Setting tidak ditemukan untuk Telegram", zap.Error(err))
		return
	}

	if !utils.IsTelegramEnabled(setting) {
		logger.Log.Warn("Telegram config incomplete")
		return
	}

	// Gunakan data booked yang sudah di-passing, bukan query lagi
	message := s.formatTelegramMessage(booked, setting)
	if message == "" {
		return
	}

	cfg := utils.TelegramConfig{
		BotToken: *setting.TelegramBotToken,
		ChatID:   *setting.TelegramChatId,
	}

	if err := utils.SendTelegramHTML(cfg, message); err != nil {
		logger.Log.Error("Gagal mengirim notifikasi Telegram ke API", zap.Error(err))
		return
	}

	logger.Log.Info("Notifikasi Telegram berhasil dikirim")
}

func (s *BookedService) formatTelegramMessage(booked *models.Booked, setting *models.Setting) string {
	// Pastikan User ada, karena ini data wajib
	if booked.User.ID == "" {
		return ""
	}

	var (
		className, isGroupStr, locationStr, priceStr string
		dateStr, timeStr                             string
		targetTime                                   time.Time
	)

	// 1. Tentukan Timezone
	timezone := "UTC"
	if setting.Timezone != nil && *setting.Timezone != "" {
		timezone = *setting.Timezone
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	// 2. Ekstraksi Data berdasarkan Tipe Booking (Schedule vs Event)
	if booked.ScheduleId != nil && booked.Schedule != nil {
		// LOGIKA UNTUK SCHEDULE RUTIN
		category := booked.Schedule.Categories
		className = category.Name
		isGroupStr = "No"
		if category.IsGroup {
			isGroupStr = "Yes"
		}
		if category.Location != nil {
			locationStr = *category.Location
		}
		priceStr = s.getPrice(category.IsFree, category.IsPayAsYouWish, category.Price)
		targetTime = booked.Schedule.DateTime

	} else if booked.EventId != nil && booked.Event != nil {
		// LOGIKA UNTUK EVENT SEKALI JALAN
		event := booked.Event
		className = "[EVENT] " + event.Name
		isGroupStr = "Yes (Event)"        // Event biasanya grup/massal
		locationStr = "Check Description" // Sesuai model Event Anda (karena model Event belum ada field location)

		// Event harganya simpel (float64)
		if event.Price == 0 {
			priceStr = "Free"
		} else {
			priceStr = fmt.Sprintf("%.2f", event.Price)
		}
		targetTime = event.StartAt
	} else {
		// Jika data tidak lengkap
		return ""
	}

	// 3. Format Waktu
	scheduleDate := targetTime.In(loc)
	dateStr = scheduleDate.Format("02 January 2006")
	timeStr = scheduleDate.Format("15:04")

	// 4. Format Sosial Media User
	user := booked.User
	socialMedia := ""
	if user.Ig != nil && *user.Ig != "" {
		socialMedia = *user.Ig
	}
	if user.Fb != nil && *user.Fb != "" {
		if socialMedia != "" {
			socialMedia += " / "
		}
		socialMedia += *user.Fb
	}
	if socialMedia == "" {
		socialMedia = "N/A"
	}

	// 5. Build Message
	message := fmt.Sprintf(
		"<b>New Booking Request</b>\n"+
			"<b>Class:</b> %s\n"+
			"<b>is Group:</b> %s\n"+
			"<b>Date:</b> %s\n"+
			"<b>Time:</b> %s\n"+
			"<b>Location:</b> %s\n"+
			"<b>Price:</b> %s\n"+
			"<b>Name:</b> %s\n"+
			"<b>Phone:</b> %s\n"+
			"<b>Email:</b> %s\n"+
			"<b>Sosial Media:</b> %s",
		className,
		isGroupStr,
		dateStr,
		timeStr,
		locationStr,
		priceStr,
		user.FullName,
		user.Phone,
		user.Email,
		socialMedia,
	)

	return message
}

func (s *BookedService) getPrice(isFree, isPayAsYouWish bool, price *float64) string {
	if isFree {
		return "Free"
	}
	if isPayAsYouWish {
		return "Pay as you wish"
	}
	if price != nil {
		return fmt.Sprintf("Rp %.0f", *price)
	}
	return "N/A"
}

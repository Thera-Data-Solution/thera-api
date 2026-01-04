package services

import (
	"errors"
	"fmt"
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
}

func NewBookedService(bookingRepo *repositories.BookedRepository, scheduleRepo *repositories.SchedulesRepository, settingRepo repositories.SettingRepo) *BookedService {
	return &BookedService{
		BookingRepo:  bookingRepo,
		ScheduleRepo: scheduleRepo,
		SettingRepo:  settingRepo,
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
	if !schedule.Categories.IsGroup {
		schedule.Status = "BOOKED"
		if err := s.ScheduleRepo.Update(schedule); err != nil {
			return err
		}
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
	bookedFull, err := s.BookingRepo.GetLatestByUserAndSchedule(userId, scheduleId, tenantId)
	if err == nil {
		go s.sendTelegramNotificationAsync(bookedFull, tenantId)
		go s.sendDiscordNotification(userId, scheduleId, tenantId)
	} else {
		logger.Log.Error("Gagal mengambil data untuk notifikasi", zap.Error(err))
	}

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

func (s *BookedService) AddTestimoni(id string, testimoni *string, anonymous *bool, showTesti *bool, tenantId string) (*models.Booked, error) {
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
	var testimonis = ""
	if testimoni == nil || *testimoni == "" {
		testimonis = *booked.Testimoni
	} else {
		testimonis = *testimoni
	}

	booked.Testimoni = &testimonis
	booked.ShowTesti = showTesti
	booked.Anonymous = anonymous
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

// sendDiscordNotification sends a Discord notification about new booking
func (s *BookedService) sendDiscordNotification(userId, scheduleId, tenantId string) {
	// Get setting to check if Discord is enabled
	setting, err := s.SettingRepo.FindByTenantId(tenantId)
	if err != nil {
		logger.Log.Warn("Setting tidak ditemukan untuk Discord", zap.String("tenantId", tenantId), zap.Error(err))
		return
	}

	// Check if Discord is enabled
	if !utils.IsDiscordEnabled(setting) {
		logger.Log.Debug("Discord tidak diaktifkan untuk tenant", zap.String("tenantId", tenantId))
		return
	}

	booked, err := s.BookingRepo.GetLatestByUserAndSchedule(userId, scheduleId, tenantId)
	if err != nil {
		logger.Log.Error("Gagal mengambil data booking untuk Discord", zap.String("userId", userId), zap.String("scheduleId", scheduleId), zap.Error(err))
		return
	}

	message := s.formatDiscordMessage(booked, setting)
	if message == "" {
		logger.Log.Warn("Pesan Discord kosong", zap.String("bookingId", booked.ID))
		return
	}

	cfg := utils.DiscordWebhookConfig{
		WebhookURL: *setting.DiscordReportId,
	}

	if err := utils.SendDiscordWebhook(cfg, message); err != nil {
		logger.Log.Error("Gagal mengirim notifikasi Discord", zap.String("bookingId", booked.ID), zap.Error(err))
		return
	}

	logger.Log.Info("Notifikasi Discord berhasil dikirim", zap.String("bookingId", booked.ID))
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
	if booked.User.ID == "" || booked.Schedule.ID == "" {
		return ""
	}

	schedule := booked.Schedule
	category := schedule.Categories
	user := booked.User

	timezone := "UTC"
	if setting.Timezone != nil && *setting.Timezone != "" {
		timezone = *setting.Timezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Log.Warn("Gagal load timezone, menggunakan UTC", zap.String("timezone", timezone), zap.Error(err))
		loc = time.UTC
	}

	scheduleDate := schedule.DateTime.In(loc)
	dateStr := scheduleDate.Format("02 January 2006")
	timeStr := scheduleDate.Format("15:04")

	isGroupStr := "No"
	if category.IsGroup {
		isGroupStr = "Yes"
	}

	priceStr := s.getPrice(category.IsFree, category.IsPayAsYouWish, category.Price)

	locationStr := "N/A"
	if category.Location != nil && *category.Location != "" {
		locationStr = *category.Location
	}

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
		category.Name,
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

func (s *BookedService) formatDiscordMessage(booked *models.Booked, setting *models.Setting) string {
	if booked.User.ID == "" || booked.Schedule.ID == "" {
		return ""
	}

	schedule := booked.Schedule
	category := schedule.Categories
	user := booked.User

	timezone := "UTC"
	if setting.Timezone != nil && *setting.Timezone != "" {
		timezone = *setting.Timezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Log.Warn("Gagal load timezone, menggunakan UTC", zap.String("timezone", timezone), zap.Error(err))
		loc = time.UTC
	}

	scheduleDate := schedule.DateTime.In(loc)
	dateStr := scheduleDate.Format("02 January 2006")
	timeStr := scheduleDate.Format("15:04")

	isGroupStr := "No"
	if category.IsGroup {
		isGroupStr = "Yes"
	}

	priceStr := s.getPrice(category.IsFree, category.IsPayAsYouWish, category.Price)

	locationStr := "N/A"
	if category.Location != nil && *category.Location != "" {
		locationStr = *category.Location
	}

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

	message := fmt.Sprintf(
		"**New Booking Request**\n"+
			"**Class:** %s\n"+
			"**is Group:** %s\n"+
			"**Date:** %s\n"+
			"**Time:** %s\n"+
			"**Location:** %s\n"+
			"**Price:** %s\n"+
			"**Name:** %s\n"+
			"**Phone:** %s\n"+
			"**Email:** %s\n"+
			"**Sosial Media:** %s",
		category.Name,
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

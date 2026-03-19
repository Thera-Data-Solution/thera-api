package services

import (
	"errors"
	"fmt"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"thera-api/utils"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetService struct {
	ResetRepo   *repositories.PasswordResetRepository
	UserRepo    *repositories.UserRepository
	AdminRepo   *repositories.TenantUserRepository
	SettingRepo repositories.SettingRepo
}

func NewPasswordResetService(
	resetRepo *repositories.PasswordResetRepository,
	userRepo *repositories.UserRepository,
	adminRepo *repositories.TenantUserRepository,
	settingRepo repositories.SettingRepo,
) *PasswordResetService {
	return &PasswordResetService{
		ResetRepo:   resetRepo,
		UserRepo:    userRepo,
		AdminRepo:   adminRepo,
		SettingRepo: settingRepo,
	}
}

// CheckRateLimit checks if IP can make a request (3x per day, 15min cooldown)
func (s *PasswordResetService) CheckRateLimit(ipAddress string) error {
	// Check daily limit (3 requests per day)
	count, err := s.ResetRepo.CountRequestsByIP(ipAddress)
	if err != nil {
		logger.Log.Error("Gagal cek rate limit", zap.String("ip", ipAddress), zap.Error(err))
		return errors.New("gagal memvalidasi rate limit")
	}

	if count >= 3 {
		return errors.New("batas maksimal 3 permintaan per hari telah tercapai")
	}

	// Check cooldown (15 minutes)
	lastRequestTime, err := s.ResetRepo.GetLastRequestTime(ipAddress)
	if err != nil {
		logger.Log.Error("Gagal cek last request time", zap.String("ip", ipAddress), zap.Error(err))
		return errors.New("gagal memvalidasi rate limit")
	}

	if lastRequestTime != nil {
		timeSinceLastRequest := time.Since(*lastRequestTime)
		if timeSinceLastRequest < 15*time.Minute {
			remainingTime := 15*time.Minute - timeSinceLastRequest
			return fmt.Errorf("harap tunggu %d menit sebelum mengirim permintaan lagi", int(remainingTime.Minutes())+1)
		}
	}

	return nil
}

// ForgotPasswordUser handles forgot password for regular users
func (s *PasswordResetService) ForgotPasswordUser(email, tenantId, ipAddress string) error {
	logger.Log.Info("ForgotPasswordUser called", zap.String("email", email), zap.String("tenantId", tenantId))

	// Check rate limit
	if err := s.CheckRateLimit(ipAddress); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.UserRepo.FindByEmailAndTenant(email, tenantId)
	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		logger.Log.Warn("User tidak ditemukan untuk forgot password", zap.String("email", email))
		// Return success to prevent email enumeration
		return nil
	}

	// Get setting for email configuration
	setting, err := s.SettingRepo.FindByTenantId(tenantId)
	if err != nil || setting.MailSecret == nil || *setting.MailSecret == "" {
		logger.Log.Error("Setting email tidak ditemukan", zap.String("tenantId", tenantId))
		return errors.New("konfigurasi email tidak ditemukan")
	}

	// Generate reset token
	token := uuid.New().String()
	resetToken := &models.PasswordResetToken{
		Token:     token,
		Email:     email,
		IPAddress: ipAddress,
		UserType:  "user",
		TenantId:  tenantId,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		Used:      false,
	}

	if err := s.ResetRepo.Create(resetToken); err != nil {
		logger.Log.Error("Gagal membuat reset token", zap.String("email", email), zap.Error(err))
		return errors.New("gagal membuat token reset password")
	}

	// Send email
	resendClient := utils.NewResendClient(*setting.MailSecret)

	// Get app name from setting for email
	appName := "Thera"
	if setting.AppName != "" {
		appName = setting.AppName
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.getFrontendURL(tenantId), token)

	emailSubject := fmt.Sprintf("Reset Password - %s", appName)
	emailHTML := s.generateResetEmailHTML(user.FullName, resetURL, appName)

	// Use MailKey as from email, or default
	fromEmail := "noreply@resend.dev"
	if setting.MailKey != nil && *setting.MailKey != "" {
		fromEmail = *setting.MailKey
	}

	if err := resendClient.SendEmail(fromEmail, email, emailSubject, emailHTML); err != nil {
		logger.Log.Error("Gagal mengirim email reset password", zap.String("email", email), zap.Error(err))
		return errors.New("gagal mengirim email reset password")
	}

	logger.Log.Info("Email reset password berhasil dikirim", zap.String("email", email))
	return nil
}

// ForgotPasswordAdmin handles forgot password for admin users
func (s *PasswordResetService) ForgotPasswordAdmin(email, tenantId, ipAddress string) error {
	logger.Log.Info("ForgotPasswordAdmin called", zap.String("email", email), zap.String("tenantId", tenantId))

	// Check rate limit
	if err := s.CheckRateLimit(ipAddress); err != nil {
		return err
	}

	// Check if admin exists
	admin, err := s.AdminRepo.FindByEmailAndTenant(email, tenantId)
	if err != nil {
		logger.Log.Warn("Admin tidak ditemukan untuk forgot password", zap.String("email", email))
		return nil
	}

	// Get setting for email configuration
	setting, err := s.SettingRepo.FindByTenantId(tenantId)
	if err != nil || setting.MailSecret == nil || *setting.MailSecret == "" {
		logger.Log.Error("Setting email tidak ditemukan", zap.String("tenantId", tenantId))
		return errors.New("konfigurasi email tidak ditemukan")
	}

	// Generate reset token
	token := uuid.New().String()
	resetToken := &models.PasswordResetToken{
		Token:     token,
		Email:     email,
		IPAddress: ipAddress,
		UserType:  "admin",
		TenantId:  tenantId,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		Used:      false,
	}

	if err := s.ResetRepo.Create(resetToken); err != nil {
		logger.Log.Error("Gagal membuat reset token", zap.String("email", email), zap.Error(err))
		return errors.New("gagal membuat token reset password")
	}

	// Send email
	resendClient := utils.NewResendClient(*setting.MailSecret)

	// Get app name from setting for email
	appName := "Thera"
	if setting.AppName != "" {
		appName = setting.AppName
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.getFrontendURL(tenantId), token)

	emailSubject := fmt.Sprintf("Reset Password - %s", appName)
	emailHTML := s.generateResetEmailHTML(admin.FullName, resetURL, appName)

	// Use MailKey as from email, or default
	fromEmail := "noreply@resend.dev"
	if setting.MailKey != nil && *setting.MailKey != "" {
		fromEmail = *setting.MailKey
	}

	if err := resendClient.SendEmail(fromEmail, email, emailSubject, emailHTML); err != nil {
		logger.Log.Error("Gagal mengirim email reset password", zap.String("email", email), zap.Error(err))
		return errors.New("gagal mengirim email reset password")
	}

	logger.Log.Info("Email reset password berhasil dikirim", zap.String("email", email))
	return nil
}

// ResetPasswordUser resets password for regular users
func (s *PasswordResetService) ResetPasswordUser(token, newPassword string) error {
	logger.Log.Info("ResetPasswordUser called", zap.String("token", token))

	// Find and validate token
	resetToken, err := s.ResetRepo.FindByToken(token)
	if err != nil {
		logger.Log.Warn("Token reset tidak valid", zap.String("token", token))
		return errors.New("token reset password tidak valid atau sudah kedaluwarsa")
	}

	if resetToken.UserType != "user" {
		return errors.New("token tidak valid untuk user")
	}

	// Find user
	user, err := s.UserRepo.FindByEmailAndTenant(resetToken.Email, resetToken.TenantId)
	if err != nil {
		logger.Log.Error("User tidak ditemukan", zap.String("email", resetToken.Email))
		return errors.New("user tidak ditemukan")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Gagal hash password", zap.Error(err))
		return errors.New("gagal mengenkripsi password")
	}

	// Update password
	user.Password = string(hashed)
	if err := s.UserRepo.Update(user); err != nil {
		logger.Log.Error("Gagal update password", zap.String("userId", user.ID), zap.Error(err))
		return errors.New("gagal mengupdate password")
	}

	// Mark token as used
	if err := s.ResetRepo.MarkAsUsed(token); err != nil {
		logger.Log.Warn("Gagal mark token as used", zap.String("token", token), zap.Error(err))
		// Non-fatal, password already updated
	}

	logger.Log.Info("Password berhasil direset", zap.String("userId", user.ID))
	return nil
}

// ResetPasswordAdmin resets password for admin users
func (s *PasswordResetService) ResetPasswordAdmin(token, newPassword string) error {
	logger.Log.Info("ResetPasswordAdmin called", zap.String("token", token))

	// Find and validate token
	resetToken, err := s.ResetRepo.FindByToken(token)
	if err != nil {
		logger.Log.Warn("Token reset tidak valid", zap.String("token", token))
		return errors.New("token reset password tidak valid atau sudah kedaluwarsa")
	}

	if resetToken.UserType != "admin" {
		return errors.New("token tidak valid untuk admin")
	}

	// Find admin
	admin, err := s.AdminRepo.FindByEmailAndTenant(resetToken.Email, resetToken.TenantId)
	if err != nil {
		logger.Log.Error("Admin tidak ditemukan", zap.String("email", resetToken.Email))
		return errors.New("admin tidak ditemukan")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Gagal hash password", zap.Error(err))
		return errors.New("gagal mengenkripsi password")
	}

	// Update password
	admin.Password = string(hashed)
	if err := s.AdminRepo.Update(admin); err != nil {
		logger.Log.Error("Gagal update password", zap.String("adminId", admin.ID), zap.Error(err))
		return errors.New("gagal mengupdate password")
	}

	// Mark token as used
	if err := s.ResetRepo.MarkAsUsed(token); err != nil {
		logger.Log.Warn("Gagal mark token as used", zap.String("token", token), zap.Error(err))
		// Non-fatal, password already updated
	}

	logger.Log.Info("Password berhasil direset", zap.String("adminId", admin.ID))
	return nil
}

// Helper functions
func (s *PasswordResetService) getFrontendURL(tenantId string) string {
	// You can customize this based on your frontend URL configuration
	// For now, return a placeholder that should be configured per tenant
	return "https://app.theravickya.com"
}

func (s *PasswordResetService) generateResetEmailHTML(fullName, resetURL, appName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Reset Password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="background-color: #f4f4f4; padding: 20px; border-radius: 5px;">
		<h2 style="color: #333;">Reset Password</h2>
		<p>Halo %s,</p>
		<p>Kami menerima permintaan untuk mereset password akun Anda di %s.</p>
		<p>Klik tombol di bawah ini untuk mereset password Anda:</p>
		<div style="text-align: center; margin: 30px 0;">
			<a href="%s" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
		</div>
		<p>Atau salin dan tempel link berikut di browser Anda:</p>
		<p style="word-break: break-all; color: #007bff;">%s</p>
		<p><strong>Link ini akan kedaluwarsa dalam 1 jam.</strong></p>
		<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="font-size: 12px; color: #666;">Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
	</div>
</body>
</html>
`, fullName, appName, resetURL, resetURL)
}

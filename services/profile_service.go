package services

import (
	"errors"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService struct {
	UserRepo  *repositories.UserRepository
	AdminRepo *repositories.TenantUserRepository
}

func NewProfileService(userRepo *repositories.UserRepository, adminRepo *repositories.TenantUserRepository) *ProfileService {
	return &ProfileService{
		UserRepo:  userRepo,
		AdminRepo: adminRepo,
	}
}

// UpdateUserProfile updates user profile
func (s *ProfileService) UpdateUserProfile(userId string, fullName, phone, address *string, ig, fb *string, avatar *string) (*models.User, error) {
	logger.Log.Info("UpdateUserProfile called", zap.String("userId", userId))

	user, err := s.UserRepo.FindByID(userId)
	if err != nil {
		logger.Log.Warn("User tidak ditemukan", zap.String("userId", userId))
		return nil, errors.New("user tidak ditemukan")
	}

	// Update fields if provided
	if fullName != nil && *fullName != "" {
		user.FullName = *fullName
	}
	if phone != nil && *phone != "" {
		user.Phone = *phone
	}
	if address != nil {
		user.Address = address
	}
	if ig != nil {
		user.Ig = ig
	}
	if fb != nil {
		user.Fb = fb
	}
	if avatar != nil && *avatar != "" {
		user.Avatar = avatar
	}

	if err := s.UserRepo.Update(user); err != nil {
		logger.Log.Error("Gagal update user profile", zap.String("userId", userId), zap.Error(err))
		return nil, errors.New("gagal mengupdate profile")
	}

	// Clear password before returning
	user.Password = ""
	logger.Log.Info("User profile berhasil diupdate", zap.String("userId", userId))
	return user, nil
}

// UpdateAdminProfile updates admin profile
func (s *ProfileService) UpdateAdminProfile(adminId string, fullName *string, avatar *string) (*models.TenantUser, error) {
	logger.Log.Info("UpdateAdminProfile called", zap.String("adminId", adminId))

	admin, err := s.AdminRepo.FindByID(adminId)
	if err != nil {
		logger.Log.Warn("Admin tidak ditemukan", zap.String("adminId", adminId))
		return nil, errors.New("admin tidak ditemukan")
	}

	// Update fields if provided
	if fullName != nil && *fullName != "" {
		admin.FullName = *fullName
	}
	if avatar != nil && *avatar != "" {
		admin.Avatar = avatar
	}

	if err := s.AdminRepo.Update(admin); err != nil {
		logger.Log.Error("Gagal update admin profile", zap.String("adminId", adminId), zap.Error(err))
		return nil, errors.New("gagal mengupdate profile")
	}

	// Clear password before returning
	admin.Password = ""
	logger.Log.Info("Admin profile berhasil diupdate", zap.String("adminId", adminId))
	return admin, nil
}

// UpdateUserPassword updates user password
func (s *ProfileService) UpdateUserPassword(userId, oldPassword, newPassword string) error {
	logger.Log.Info("UpdateUserPassword called", zap.String("userId", userId))

	user, err := s.UserRepo.FindByID(userId)
	if err != nil {
		logger.Log.Warn("User tidak ditemukan", zap.String("userId", userId))
		return errors.New("user tidak ditemukan")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		logger.Log.Warn("Password lama salah", zap.String("userId", userId))
		return errors.New("password lama salah")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Gagal hash password", zap.Error(err))
		return errors.New("gagal mengenkripsi password")
	}

	user.Password = string(hashed)
	if err := s.UserRepo.Update(user); err != nil {
		logger.Log.Error("Gagal update password", zap.String("userId", userId), zap.Error(err))
		return errors.New("gagal mengupdate password")
	}

	logger.Log.Info("Password user berhasil diupdate", zap.String("userId", userId))
	return nil
}

// UpdateAdminPassword updates admin password
func (s *ProfileService) UpdateAdminPassword(adminId, oldPassword, newPassword string) error {
	logger.Log.Info("UpdateAdminPassword called", zap.String("adminId", adminId))

	admin, err := s.AdminRepo.FindByID(adminId)
	if err != nil {
		logger.Log.Warn("Admin tidak ditemukan", zap.String("adminId", adminId))
		return errors.New("admin tidak ditemukan")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword))
	if err != nil {
		logger.Log.Warn("Password lama salah", zap.String("adminId", adminId))
		return errors.New("password lama salah")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Gagal hash password", zap.Error(err))
		return errors.New("gagal mengenkripsi password")
	}

	admin.Password = string(hashed)
	if err := s.AdminRepo.Update(admin); err != nil {
		logger.Log.Error("Gagal update password", zap.String("adminId", adminId), zap.Error(err))
		return errors.New("gagal mengupdate password")
	}

	logger.Log.Info("Password admin berhasil diupdate", zap.String("adminId", adminId))
	return nil
}


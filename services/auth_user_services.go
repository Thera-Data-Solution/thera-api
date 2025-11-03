package services

import (
	"errors"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserService struct {
	UserRepo    *repositories.UserRepository
	SessionRepo *repositories.SessionRepository
	TenantRepo  *repositories.TenantRepository
}

func (s *AuthUserService) RegisterUser(
	email,
	password,
	fullName,
	phone,
	address,
	ig,
	fb,
	tenantId string,
) (*models.Session, error) {
	logger.Log.Info("RegisterUser called", zap.String("email", email), zap.String("tenantId", tenantId))

	_, err := s.TenantRepo.FindByID(tenantId)
	if err != nil {
		logger.Log.Warn("Tenant tidak ditemukan", zap.String("tenantId", tenantId))
		return nil, errors.New("tenant tidak ditemukan")
	}

	if fullName == "" || phone == "" || email == "" || (ig == "" && fb == "") {
		logger.Log.Warn("Validasi input gagal", zap.String("email", email))
		return nil, errors.New("data input tidak lengkap")
	}

	existing, _ := s.UserRepo.FindByEmailAndTenant(email, tenantId)
	if existing.ID != "" {
		logger.Log.Warn("Email sudah terdaftar", zap.String("email", email))
		return nil, errors.New("pengguna dengan email tersebut sudah terdaftar")
	}

	existingPhone, _ := s.UserRepo.FindByPhoneAndTenant(phone, tenantId)
	if existingPhone.ID != "" {
		logger.Log.Warn("Nomor telepon sudah terdaftar", zap.String("phone", phone))
		return nil, errors.New("pengguna dengan nomor handphone tersebut sudah terdaftar")
	}

	avatar := "https://api.dicebear.com/9.x/fun-emoji/svg?seed=" + email
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Gagal hash password", zap.Error(err))
		return nil, errors.New("gagal mengenkripsi password")
	}

	user := models.User{
		Email:    email,
		Password: string(hashed),
		FullName: fullName,
		Avatar:   &avatar,
		Fb:       &fb,
		Ig:       &ig,
		Address:  &address,
		TenantId: tenantId,
	}

	if err := s.UserRepo.Create(&user); err != nil {
		logger.Log.Error("Gagal membuat user", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	token := uuid.NewString()
	session := &models.Session{
		Token:     token,
		UserId:    &user.ID,
		TenantId:  tenantId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		logger.Log.Error("Gagal membuat session", zap.String("userId", user.ID), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User registered successfully", zap.String("userId", user.ID))
	return session, nil
}

func (s *AuthUserService) LoginUser(email, password, tenantId string) (*models.Session, error) {
	logger.Log.Info("LoginUser called", zap.String("email", email), zap.String("tenantId", tenantId))

	user, err := s.UserRepo.FindByEmailAndTenant(email, tenantId)
	if err != nil {
		logger.Log.Warn("Pengguna tidak ditemukan", zap.String("email", email))
		return nil, errors.New("pengguna tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Log.Warn("Password salah", zap.String("email", email))
		return nil, errors.New("password salah")
	}

	session := &models.Session{
		Token:     uuid.New().String(),
		UserId:    &user.ID,
		TenantId:  user.TenantId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.SessionRepo.DeleteByTenantUserId(session); err != nil {
		logger.Log.Error("Gagal hapus session lama", zap.String("userId", user.ID), zap.Error(err))
		return nil, err
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		logger.Log.Error("Gagal buat session baru", zap.String("userId", user.ID), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User login berhasil", zap.String("userId", user.ID))
	return session, nil
}

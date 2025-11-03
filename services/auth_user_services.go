package services

import (
	"errors"
	"fmt"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"github.com/google/uuid"
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
	_, err := s.TenantRepo.FindByID(tenantId)
	if err != nil {
		return nil, errors.New("tenant tidak ditemukan")
	}
	if fullName == "" {
		return nil, errors.New("nama lengkap wajib diisi")
	}
	if phone == "" {
		return nil, errors.New("nomor telepon wajib diisi")
	}
	if email == "" {
		return nil, errors.New("email wajib diisi")
	}
	if ig == "" && fb == "" {
		return nil, errors.New("minimal salah satu dari Instagram atau Facebook wajib diisi")
	}
	fmt.Println(password)
	existing, _ := s.UserRepo.FindByEmailAndTenant(email, tenantId)
	if existing.ID != "" {
		return nil, errors.New("maaf, pengguna dengan email tersebut sudah terdaftar")
	}

	existingPhone, _ := s.UserRepo.FindByPhoneAndTenant(phone, tenantId)
	if existingPhone.ID != "" {
		return nil, errors.New("maaf, pengguna dengan nomor handphone tersebut sudah terdaftar")
	}

	// buat avatar random
	avatar := "https://api.dicebear.com/9.x/fun-emoji/svg?seed=" + email

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
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

	// simpan user baru
	if err := s.UserRepo.Create(&user); err != nil {
		return nil, err
	}

	// otomatis login => buat token session
	token := uuid.NewString()
	session := &models.Session{
		Token:     token,
		UserId:    &user.ID,
		TenantId:  tenantId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 hari
		CreatedAt: time.Now(),
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthUserService) LoginUser(email, password, tenantId string) (*models.Session, error) {
	user, err := s.UserRepo.FindByEmailAndTenant(email, tenantId)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
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
		return nil, err
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

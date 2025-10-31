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

type AuthAdminService struct {
	AdminRepo   *repositories.TenantUserRepository
	SessionRepo *repositories.SessionRepository
	TenantRepo  *repositories.TenantRepository
}

func (s *AuthAdminService) RegisterAdmin(email, password, fullName, tenantId string) (*models.Session, error) {
	existing, _ := s.AdminRepo.FindByEmailAndTenant(email, tenantId)
	if existing.ID != "" {
		return nil, errors.New("pengguna dengan email tersebut sudah terdaftar")
	}

	if fullName == "" {
		return nil, errors.New("fullname wajib diisi")
	}

	// buat avatar random
	avatar := "https://api.dicebear.com/9.x/fun-emoji/svg?seed=" + email

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	admin := models.TenantUser{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashed),
		FullName: fullName,
		Avatar:   &avatar,
		Role:     "USER",
		TenantId: tenantId,
	}

	// simpan user baru
	if err := s.AdminRepo.Create(&admin); err != nil {
		return nil, err
	}

	// otomatis login => buat token session
	token := uuid.NewString()
	session := &models.Session{
		ID:           uuid.New().String(),
		Token:        token,
		TenantUserId: &admin.ID,
		TenantId:     tenantId,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 hari
		CreatedAt:    time.Now(),
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthAdminService) LoginAdmin(email, password, tenantId string) (*models.Session, error) {
	admin, err := s.AdminRepo.FindByEmailAndTenant(email, tenantId)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	fmt.Println("HASH:", admin.Password)
	fmt.Println("RAW :", password)
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		fmt.Println("bcrypt error:", err)
		return nil, errors.New("password salah")
	}

	session := &models.Session{
		ID:           uuid.New().String(),
		Token:        uuid.New().String(),
		TenantUserId: &admin.ID,
		TenantId:     admin.TenantId,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	if err := s.SessionRepo.DeleteByTenantUserId(session); err != nil {
		return nil, err
	}

	if err := s.SessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

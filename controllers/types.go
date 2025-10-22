package controllers

import "thera-api/repository"

type AuthController struct {
	UserRepo       *repository.UserRepository
	TenantUserRepo *repository.TenantUserRepository
	SessionRepo    *repository.SessionRepository
	TenantRepo     *repository.TenantRepo
}

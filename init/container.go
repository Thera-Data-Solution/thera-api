package init

import (
	"thera-api/config"
	handlers "thera-api/handles"
	"thera-api/middlewares"
	"thera-api/repositories"
	"thera-api/services"
)

type Container struct {
	UserHandler   *handlers.AuthUserHandler
	AdminHandler  *handlers.AuthAdminHandler
	Middlewares   *middlewares.IsAuthMiddleware
	AtLeastAdmin  *middlewares.IsAdminMiddleware
	OnlySU        *middlewares.IsSUMiddleware
	TenantHandler *handlers.TenantHandler
}

func NewContainer() *Container {
	// koneksi DB (pastikan sudah connect di main.go)
	db := config.DB

	userRepo := &repositories.UserRepository{DB: db}
	adminRepo := &repositories.TenantUserRepository{DB: db}
	sessionRepo := &repositories.SessionRepository{DB: db}

	tenantRepo := &repositories.TenantRepository{DB: db}

	authUserService := &services.AuthUserService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
	authAdminService := &services.AuthAdminService{
		AdminRepo:   adminRepo,
		SessionRepo: sessionRepo,
	}

	tenantService := &services.TenantService{
		TenantRepo: tenantRepo,
	}

	userHandler := &handlers.AuthUserHandler{Service: authUserService}
	adminHandler := &handlers.AuthAdminHandler{Service: authAdminService}
	tenantHandler := &handlers.TenantHandler{Service: tenantService}

	authAdminMiddleware := &middlewares.IsAuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		AdminRepo:   adminRepo,
	}
	atLeastAdminMiddleware := &middlewares.IsAdminMiddleware{}
	onlySUMiddleware := &middlewares.IsSUMiddleware{}

	return &Container{
		UserHandler:   userHandler,
		AdminHandler:  adminHandler,
		Middlewares:   authAdminMiddleware,
		AtLeastAdmin:  atLeastAdminMiddleware,
		TenantHandler: tenantHandler,
		OnlySU:        onlySUMiddleware,
	}
}

package init

import (
	"thera-api/config"
	handlers "thera-api/handles"
	"thera-api/middlewares"
	"thera-api/repositories"
	"thera-api/services"
)

type Container struct {
	UserHandler  *handlers.AuthUserHandler
	AdminHandler *handlers.AuthAdminHandler
	Middlewares  *middlewares.IsAuthMiddleware
}

func NewContainer() *Container {
	// koneksi DB (pastikan sudah connect di main.go)
	db := config.DB

	userRepo := &repositories.UserRepository{DB: db}
	adminRepo := &repositories.TenantUserRepository{DB: db}
	sessionRepo := &repositories.SessionRepository{DB: db}

	authUserService := &services.AuthUserService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
	authAdminService := &services.AuthAdminService{
		AdminRepo:   adminRepo,
		SessionRepo: sessionRepo,
	}

	userHandler := &handlers.AuthUserHandler{Service: authUserService}
	adminHandler := &handlers.AuthAdminHandler{Service: authAdminService}
	authAdminMiddleware := &middlewares.IsAuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		AdminRepo:   adminRepo,
	}

	return &Container{
		UserHandler:  userHandler,
		AdminHandler: adminHandler,
		Middlewares:  authAdminMiddleware,
	}
}

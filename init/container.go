package init

import (
	"thera-api/config"
	"thera-api/handlers"
	"thera-api/middlewares"
	"thera-api/repositories"
	"thera-api/services"
)

type Container struct {
	UserHandler        *handlers.AuthUserHandler
	AdminHandler       *handlers.AuthAdminHandler
	Middlewares        *middlewares.IsAuthMiddleware
	AtLeastAdmin       *middlewares.IsAdminMiddleware
	OnlySU             *middlewares.IsSUMiddleware
	TenantHandler      *handlers.TenantHandler
	CategoryHandler    *handlers.CategoriesHandler
	ScheduleHandler    *handlers.SchedulesHandler
	BookHandler        *handlers.BookedHandler
	HeroHandler        *handlers.HeroHandler
	LinkHandler        *handlers.LinkHandler
	ArticleHandler     *handlers.ArticleHandler
	GalleryHandler     *handlers.GalleryHandler
	TranslationHandler *handlers.TranslationHandler
	SettingHandler     *handlers.SettingHandler
}

func NewContainer() *Container {
	db := config.DB

	userRepo := &repositories.UserRepository{DB: db}
	adminRepo := &repositories.TenantUserRepository{DB: db}
	sessionRepo := &repositories.SessionRepository{DB: db}

	tenantRepo := &repositories.TenantRepository{DB: db}
	categoriesRepo := &repositories.CategoriesRepository{DB: db}
	scheduleRepo := &repositories.SchedulesRepository{DB: db}
	bookingRepo := &repositories.BookedRepository{DB: db}
	heroRepo := &repositories.HeroRepository{DB: db}
	linkRepo := &repositories.LinkRepository{DB: db}
	articleRepo := &repositories.ArticleRepository{DB: db}
	galleryRepo := &repositories.GalleryRepository{DB: db}
	translateRepo := &repositories.TranslationRepository{DB: db}
	settingRepo := repositories.NewSettingRepo(db)

	authUserService := &services.AuthUserService{UserRepo: userRepo, SessionRepo: sessionRepo, TenantRepo: tenantRepo}
	authAdminService := &services.AuthAdminService{AdminRepo: adminRepo, SessionRepo: sessionRepo, TenantRepo: tenantRepo}
	tenantService := &services.TenantService{TenantRepo: tenantRepo}
	categoryService := &services.CategoriesService{CategoriesRepo: categoriesRepo}
	scheduleService := &services.SchedulesService{SchedulesRepo: scheduleRepo}
	bookingService := &services.BookedService{BookingRepo: bookingRepo, ScheduleRepo: scheduleRepo}
	heroService := &services.HeroService{Repo: heroRepo}
	linkService := &services.LinkService{Repo: linkRepo}
	articleService := &services.ArticleService{Repo: articleRepo}
	galleryService := &services.GalleryService{GalleryRepo: galleryRepo}
	translateService := &services.TranslationService{Repo: translateRepo}
	settingService := services.NewSettingService(settingRepo)

	userHandler := &handlers.AuthUserHandler{Service: authUserService}
	adminHandler := &handlers.AuthAdminHandler{Service: authAdminService}
	tenantHandler := &handlers.TenantHandler{Service: tenantService}
	categoryHandler := &handlers.CategoriesHandler{Service: categoryService}
	scheduleHandler := &handlers.SchedulesHandler{Service: scheduleService}
	bookHandler := &handlers.BookedHandler{Service: bookingService}
	heroHandler := &handlers.HeroHandler{Service: heroService}
	linkHandler := &handlers.LinkHandler{Service: linkService}
	articleHandler := &handlers.ArticleHandler{Service: articleService}
	galleryHandler := &handlers.GalleryHandler{Service: galleryService}
	translateHandler := &handlers.TranslationHandler{Service: translateService}
	settingHandler := handlers.NewSettingHandler(settingService)

	authAdminMiddleware := &middlewares.IsAuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		AdminRepo:   adminRepo,
	}
	atLeastAdminMiddleware := &middlewares.IsAdminMiddleware{}
	onlySUMiddleware := &middlewares.IsSUMiddleware{}

	return &Container{
		UserHandler:        userHandler,
		AdminHandler:       adminHandler,
		Middlewares:        authAdminMiddleware,
		AtLeastAdmin:       atLeastAdminMiddleware,
		TenantHandler:      tenantHandler,
		OnlySU:             onlySUMiddleware,
		CategoryHandler:    categoryHandler,
		ScheduleHandler:    scheduleHandler,
		BookHandler:        bookHandler,
		HeroHandler:        heroHandler,
		LinkHandler:        linkHandler,
		ArticleHandler:     articleHandler,
		GalleryHandler:     galleryHandler,
		TranslationHandler: translateHandler,
		SettingHandler:     settingHandler,
	}
}

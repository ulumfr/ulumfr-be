package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ulumfr/ulumfr-be/internal/config"
	"github.com/ulumfr/ulumfr-be/internal/handler"
	"github.com/ulumfr/ulumfr-be/internal/middleware"
	"github.com/ulumfr/ulumfr-be/internal/repository"
	"github.com/ulumfr/ulumfr-be/internal/service"
	"github.com/ulumfr/ulumfr-be/internal/storage"
	"github.com/ulumfr/ulumfr-be/prisma/db"

	_ "github.com/ulumfr/ulumfr-be/docs" // swagger docs
)

func main() {
	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("APP_ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().
		Str("env", cfg.AppEnv).
		Str("port", cfg.Port).
		Msg("Starting server...")

	// Initialize Prisma client
	prismaClient := db.NewClient()
	if err := prismaClient.Prisma.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		if err := prismaClient.Prisma.Disconnect(); err != nil {
			log.Error().Err(err).Msg("Failed to disconnect from database")
		}
	}()

	log.Info().Msg("Connected to database")

	// Initialize R2 storage client
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to initialize R2 client - file uploads will be disabled")
	}

	// Initialize repositories
	repos := repository.NewRepositories(prismaClient)

	// Initialize services
	services := service.NewServices(repos, r2Client, cfg)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "ulumfr-be",
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           120 * time.Second,
		DisableStartupMessage: cfg.IsProduction(),
		ErrorHandler:          handler.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinOrigins(cfg.AllowedOrigins),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(repos.User, cfg)

	// Setup routes
	setupRoutes(app, services, authMiddleware, cfg)

	// Start server in goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}

func setupRoutes(app *fiber.App, services *service.Services, authMiddleware *middleware.AuthMiddleware, cfg *config.Config) {
	// Health check
	app.Get("/health", handler.HealthCheck)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1 group
	api := app.Group("/api/v1")

	// ===========================================
	// Auth routes (no authentication required)
	// ===========================================
	auth := api.Group("/auth")
	{
		// Rate limit for auth endpoints
		authLimiter := limiter.New(limiter.Config{
			Max:        10,
			Expiration: time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests. Please try again later.",
				})
			},
		})

		auth.Post("/register", authLimiter, services.Auth.Register)
		auth.Post("/login", authLimiter, services.Auth.Login)
		auth.Post("/refresh", services.Auth.RefreshToken)
		auth.Get("/me", authMiddleware.RequireAuth(), services.Auth.Me)
		auth.Post("/logout", authMiddleware.RequireAuth(), services.Auth.Logout)
		auth.Post("/logout-all", authMiddleware.RequireAuth(), services.Auth.LogoutAll)
	}

	// ===========================================
	// Public routes (no authentication required)
	// ===========================================
	public := api.Group("/public")
	{
		// Projects
		public.Get("/projects", services.Project.List)
		public.Get("/projects/:slug", services.Project.GetBySlug)

		// Categories & Tags
		public.Get("/categories", services.Category.List)
		public.Get("/tags", services.Tag.List)

		// Career & Education
		public.Get("/careers", services.Career.List)
		public.Get("/educations", services.Education.List)

		// Resume
		public.Get("/resume", services.Resume.GetActive)

		// Contact (with rate limiting)
		contactLimiter := limiter.New(limiter.Config{
			Max:        5,
			Expiration: time.Duration(cfg.RateLimitWindowSeconds) * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests. Please try again later.",
				})
			},
		})
		public.Post("/contact", contactLimiter, services.Contact.Create)
	}

	// ===========================================
	// Admin routes (JWT + Admin role required)
	// ===========================================
	admin := api.Group("/admin", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin())
	{
		// Projects
		admin.Get("/projects", services.Project.AdminList)
		admin.Post("/projects", services.Project.Create)
		admin.Get("/projects/:id", services.Project.GetByID)
		admin.Put("/projects/:id", services.Project.Update)
		admin.Delete("/projects/:id", services.Project.Delete)

		// Categories
		admin.Get("/categories", services.Category.AdminList)
		admin.Post("/categories", services.Category.Create)
		admin.Put("/categories/:id", services.Category.Update)
		admin.Delete("/categories/:id", services.Category.Delete)

		// Tags
		admin.Get("/tags", services.Tag.AdminList)
		admin.Post("/tags", services.Tag.Create)
		admin.Put("/tags/:id", services.Tag.Update)
		admin.Delete("/tags/:id", services.Tag.Delete)

		// Careers
		admin.Get("/careers", services.Career.AdminList)
		admin.Post("/careers", services.Career.Create)
		admin.Put("/careers/:id", services.Career.Update)
		admin.Delete("/careers/:id", services.Career.Delete)

		// Educations
		admin.Get("/educations", services.Education.AdminList)
		admin.Post("/educations", services.Education.Create)
		admin.Put("/educations/:id", services.Education.Update)
		admin.Delete("/educations/:id", services.Education.Delete)

		// Resumes
		admin.Get("/resumes", services.Resume.List)
		admin.Post("/resumes", services.Resume.Create)
		admin.Put("/resumes/:id", services.Resume.Update)
		admin.Delete("/resumes/:id", services.Resume.Delete)
		admin.Post("/resumes/:id/activate", services.Resume.Activate)

		// Contacts
		admin.Get("/contacts", services.Contact.List)
		admin.Get("/contacts/:id", services.Contact.GetByID)
		admin.Put("/contacts/:id/read", services.Contact.MarkAsRead)
		admin.Delete("/contacts/:id", services.Contact.Delete)

		// File Upload
		admin.Post("/upload-url", services.Upload.GetPresignedURL)
	}
}

func joinOrigins(origins []string) string {
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}

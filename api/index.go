package handler

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"github.com/ulumfr/ulumfr-be/pkg/config"
	"github.com/ulumfr/ulumfr-be/pkg/handler"
	"github.com/ulumfr/ulumfr-be/pkg/middleware"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/service"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
	"github.com/ulumfr/ulumfr-be/prisma/db"

	_ "github.com/ulumfr/ulumfr-be/docs"
)

var (
	app  *fiber.App
	once sync.Once
)

// Handler is the entrypoint for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	adaptor.FiberApp(app)(w, r)
}

func initApp() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize Prisma client
	prismaClient := db.NewClient()
	if err := prismaClient.Prisma.Connect(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Initialize R2 storage client
	r2Client, _ := storage.NewR2Client(cfg)

	// Initialize repositories
	repos := repository.NewRepositories(prismaClient)

	// Initialize services
	services := service.NewServices(repos, r2Client, cfg)

	// Initialize Fiber app
	app = fiber.New(fiber.Config{
		AppName:               "ulumfr-be",
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		DisableStartupMessage: true,
		ErrorHandler:          handler.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
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
}

func setupRoutes(app *fiber.App, services *service.Services, authMiddleware *middleware.AuthMiddleware, cfg *config.Config) {
	// Health check
	app.Get("/health", handler.HealthCheck)

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1
	api := app.Group("/v1")

	// Auth routes
	auth := api.Group("/auth")
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
	auth.Post("/register", authLimiter, services.Auth.Register)
	auth.Post("/login", authLimiter, services.Auth.Login)
	auth.Post("/refresh", services.Auth.RefreshToken)
	auth.Get("/me", authMiddleware.RequireAuth(), services.Auth.Me)
	auth.Post("/logout", authMiddleware.RequireAuth(), services.Auth.Logout)
	auth.Post("/logout-all", authMiddleware.RequireAuth(), services.Auth.LogoutAll)
	auth.Put("/profile", authMiddleware.RequireAuth(), services.Auth.UpdateProfile)

	// Public routes
	public := api.Group("/public")
	public.Get("/projects", services.Project.List)
	public.Get("/projects/:slug", services.Project.GetBySlug)
	public.Get("/categories", services.Category.List)
	public.Get("/tags", services.Tag.List)
	public.Get("/careers", services.Career.List)
	public.Get("/educations", services.Education.List)
	public.Get("/resume", services.Resume.GetActive)
	contactLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Duration(cfg.RateLimitWindowSeconds) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
	public.Post("/contact", contactLimiter, services.Contact.Create)

	// Admin routes
	admin := api.Group("/admin", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin())
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
	// Upload
	admin.Post("/upload-url", services.Upload.GetPresignedURL)
	// Users
	admin.Get("/users", services.Auth.ListUsers)
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

// For local development
func init() {
	if os.Getenv("VERCEL") == "" {
		// Not running on Vercel, skip initialization
		return
	}
}

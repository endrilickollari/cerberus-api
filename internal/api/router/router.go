package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "remote-server-api/docs" // Import for swagger docs
	"remote-server-api/internal/api/handlers"
	"remote-server-api/internal/domain/auth"
	"remote-server-api/internal/domain/docker"
	"remote-server-api/internal/domain/server"
)

// New creates and configures a router with all application routes
func New(
	authService auth.Service,
	serverService server.Service,
	dockerService docker.Service,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create handlers with injected services
	authHandler := handlers.NewAuthHandler(authService)
	serverHandler := handlers.NewServerHandler(serverService)
	dockerHandler := handlers.NewDockerHandler(dockerService)

	// Authentication middleware
	authMiddleware := handlers.NewAuthMiddleware(authService)

	// Swagger documentation
	// This serves the Swagger UI at /swagger/index.html
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// Server details routes
		r.Route("/server-details", func(r chi.Router) {
			r.Get("/", serverHandler.GetBasicDetails)
			r.Get("/cpu-info", serverHandler.GetCPUInfo)
			r.Get("/disk-usage", serverHandler.GetDiskUsage)
			r.Get("/running-processes", serverHandler.GetRunningProcesses)
		})

		// Docker routes
		r.Route("/docker", func(r chi.Router) {
			r.Get("/container-details", dockerHandler.GetContainerInfo)
		})
	})

	return r
}

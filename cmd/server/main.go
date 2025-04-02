package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	dockerDomain "remote-server-api/internal/domain/docker"
	"syscall"
	"time"

	"remote-server-api/config"
	"remote-server-api/internal/api/router"
	"remote-server-api/internal/api/server"
	"remote-server-api/internal/domain/auth"
	serverDomain "remote-server-api/internal/domain/server"
	"remote-server-api/internal/infrastructure/persistence/memory"
	"remote-server-api/internal/infrastructure/ssh"
	"remote-server-api/internal/infrastructure/token"

	// Import for swagger docs
	_ "remote-server-api/docs"
)

// @title Cerberus API
// @version 2.0.0
// @description API for Cerberus - Remote Server Management
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Setup dependencies
	tokenService := token.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	sshClient := ssh.NewClient()

	// Setup repositories
	sessionRepo := memory.NewSessionRepository()

	// Setup services
	authService := auth.NewService(sessionRepo, sshClient, tokenService)
	serverService := serverDomain.NewService(sessionRepo)
	dockerService := dockerDomain.NewService(sessionRepo)

	// Setup router with all dependencies
	r := router.New(authService, serverService, dockerService)

	// Initialize HTTP server
	srv := server.NewServer(r, cfg.Server)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on :%s", cfg.Server.Port)
		log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

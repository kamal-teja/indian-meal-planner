package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meal-planner-backend/internal/api"
	"meal-planner-backend/internal/config"
	"meal-planner-backend/internal/database"
	"meal-planner-backend/internal/repository"
	"meal-planner-backend/internal/service"
	"meal-planner-backend/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel, cfg.LogFormat)

	// Initialize database
	db, err := database.Connect(cfg.MongoURI, cfg.DatabaseConfig, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Disconnect()

	// Seed database with default dishes
	if err := database.SeedDefaultDishes(db.GetDB(), logger); err != nil {
		logger.Warn("Failed to seed database", "error", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db.GetDB())

	// Initialize services
	services := service.NewServices(repos, cfg, logger)

	// Initialize API router
	router := api.NewRouter(services, cfg, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		logger.Info("API Health Check", "url", "http://localhost:"+cfg.Port+"/api/health")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server shutdown complete")
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/ryangel/ryangel-backend/internal/config"
	"github.com/ryangel/ryangel-backend/internal/database"
	"github.com/ryangel/ryangel-backend/internal/logger"
	"github.com/ryangel/ryangel-backend/internal/repository"
	"github.com/ryangel/ryangel-backend/internal/server"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
	ebuysvc "github.com/ryangel/ryangel-backend/internal/services"
)

func main() {
	_ = godotenv.Load(".env")

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	appLogger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer appLogger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		appLogger.Fatal("database connection error", zap.Error(err))
	}
	defer pool.Close()

	adminRepo := repository.NewAdminRepository(pool)
	clientRepo := repository.NewClientRepository(pool)
	authService := authsvc.NewService(adminRepo, clientRepo, cfg)
	ebuyService := ebuysvc.NewEbuyService(pool)

	srv := server.New(server.Options{
		Config:      cfg,
		DB:          pool,
		Logger:      appLogger,
		AuthService: authService,
		EbuyService: ebuyService,
	})

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server stopped", zap.Error(err))
		}
	}()

	appLogger.Info("server listening", zap.String("addr", cfg.HTTPAddr()))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful shutdown failed", zap.Error(err))
	} else {
		appLogger.Info("server shutdown complete")
	}
}

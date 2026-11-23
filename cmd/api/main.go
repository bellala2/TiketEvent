package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/bellala2/TiketEvent/internal/clock"
	"github.com/bellala2/TiketEvent/internal/config"
	"github.com/bellala2/TiketEvent/internal/database"
	"github.com/bellala2/TiketEvent/internal/handler"
	"github.com/bellala2/TiketEvent/internal/middleware"
	"github.com/bellala2/TiketEvent/internal/repository"
	"github.com/bellala2/TiketEvent/internal/response"
	"github.com/bellala2/TiketEvent/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("aplikasi berhenti karena error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// ---- Wiring: repository -> service -> handler ----
	clk := clock.System{}
	tokens := service.NewTokenManager(cfg.JWTSecret, time.Hour, clk)

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, tokens)

	  
	if err := authSvc.EnsureAdmin(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		return fmt.Errorf("gagal seed admin: %w", err)
	}

	healthHandler := handler.NewHealthHandler(db, cfg.ImportQueueSize)
	authHandler := handler.NewAuthHandler(authSvc)

	// ---- Routing ----
	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(w, service.NewError(service.CodeNotFound, "Endpoint tidak ditemukan"))
	})

	r.Get("/health", healthHandler.Check)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.With(middleware.Authenticate(tokens)).Get("/me", authHandler.Me)
		})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		fmt.Println("Server berjalan di port " + cfg.Port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		fmt.Println("Server dimatikan...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
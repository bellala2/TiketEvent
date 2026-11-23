package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menampung semua pengaturan aplikasi yang dibaca dari environment variable.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AdminEmail      string
	AdminPassword   string
	ImportWorkers   int
	ImportQueueSize int
}

// Load membaca file .env (kalau ada) lalu mengisi Config.
func Load() (*Config, error) {
	// Error diabaikan: di server asli biasanya tidak ada file .env,
	// nilai langsung datang dari environment.
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AdminEmail:      os.Getenv("ADMIN_EMAIL"),
		AdminPassword:   os.Getenv("ADMIN_PASSWORD"),
		ImportWorkers:   getEnvInt("IMPORT_WORKERS", 4),
		ImportQueueSize: getEnvInt("IMPORT_QUEUE_SIZE", 10),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL wajib diisi")
	}
	if len(cfg.JWTSecret) < 16 {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi, minimal 16 karakter")
	}
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("ADMIN_EMAIL dan ADMIN_PASSWORD wajib diisi (untuk seed akun admin)")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
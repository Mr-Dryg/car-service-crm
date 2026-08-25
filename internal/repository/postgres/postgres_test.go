package postgres

import (
	"context"
	"testing"

	"github.com/Mr-Dryg/car-service-crm/internal/config"
)

func TestNewPool(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "db_user")
	t.Setenv("DB_PASSWORD", "db_password")
	t.Setenv("DB_NAME", "car_service")
	t.Setenv("DB_SSLMODE", "disable")

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("load end with error: %v", err)
	}

	ctx := context.Background()
	pool, err := NewPool(ctx, cfg.DatabaseURL)

	if err != nil {
		t.Fatalf("new pool with error: %v", err)
	}

	defer pool.Close()
}

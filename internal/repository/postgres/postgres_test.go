package postgres

import (
	"context"
	"testing"

	"github.com/Mr-Dryg/car-service-crm/internal/config"
)

func TestNewPool(t *testing.T) {
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

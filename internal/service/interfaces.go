package service

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

type BranchRepository interface {
	Create(ctx context.Context, branch *domain.Branch) error
	GetAll(ctx context.Context) ([]domain.Branch, error)
}

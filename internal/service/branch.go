package service

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

type BranchService struct {
	branchRepo BranchRepository
}

func NewBranchRepository(repo BranchRepository) *BranchService {
	return &BranchService{
		branchRepo: repo,
	}
}

func (s *BranchService) Create(ctx context.Context, branch *domain.Branch) error {
	return s.branchRepo.Create(ctx, branch)
}

func (s *BranchService) GetAll(ctx context.Context) ([]domain.Branch, error) {
	return s.branchRepo.GetAll(ctx)
}

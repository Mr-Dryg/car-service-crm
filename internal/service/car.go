package service

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

type CarRepository interface {
	Create(ctx context.Context, car *domain.Car) error
	GetByLicensePlate(ctx context.Context, plate string) (*domain.Car, error)
	GetByUserID(ctx context.Context, userID int64) ([]domain.Car, error)
}

type CarService struct {
	carRepo CarRepository
}

func NewCarService(repo CarRepository) *CarService {
	return &CarService{
		carRepo: repo,
	}
}

func (s *CarService) Create(ctx context.Context, car *domain.Car) error {
	return s.carRepo.Create(ctx, car)
}

func (s *CarService) GetByLicensePlate(ctx context.Context, plate string) (*domain.Car, error) {
	return s.carRepo.GetByLicensePlate(ctx, plate)
}

func (s *CarService) GetByUserID(ctx context.Context, userID int64) ([]domain.Car, error) {
	return s.carRepo.GetByUserID(ctx, userID)
}

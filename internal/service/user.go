package service

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByTelegramID(ctx context.Context, tgID int64) (*domain.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User, password string) error {
	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(passwordHash)
	}
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return s.userRepo.GetByPhone(ctx, phone)
}

func (s *UserService) GetByTelegramID(ctx context.Context, tgID int64) (*domain.User, error) {
	return s.userRepo.GetByTelegramID(ctx, tgID)
}

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

type CreateUserDTO struct {
	Name       string
	Phone      string
	Email      string
	Password   string
	Role       string
	TelegramID int64
	BranchID   int64
}

func (s *UserService) Create(ctx context.Context, userDTO *CreateUserDTO) (*domain.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         userDTO.Name,
		Phone:        userDTO.Phone,
		Email:        userDTO.Email,
		PasswordHash: string(passwordHash),
		Role:         userDTO.Role,
		TelegramID:   userDTO.TelegramID,
		BranchID:     userDTO.BranchID,
	}
	return user, s.userRepo.Create(ctx, user)
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

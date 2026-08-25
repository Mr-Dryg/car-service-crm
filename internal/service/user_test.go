package service

import (
    "context"
    "testing"

    "github.com/Mr-Dryg/car-service-crm/internal/domain"
    "golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
    created *domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, user *domain.User) error {
    f.created = user
    return nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    return nil, nil
}

func (f *fakeUserRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
    return nil, nil
}

func (f *fakeUserRepo) GetByTelegramID(ctx context.Context, tgID int64) (*domain.User, error) {
    return nil, nil
}

func TestCreateHashesPassword(t *testing.T) {
    repo := &fakeUserRepo{}
    svc := NewUserService(repo)

    user := &domain.User{
        Name:     "Иван",
        Phone:    "+7 (900) 123-45-67",
        Email:    "ivan@example.com",
        Role:     "client",
    }
    password := "secret123"

    err := svc.Create(context.Background(), user, password)
    if err != nil {
        t.Fatal(err)
    }

    if user.PasswordHash == password {
        t.Error("password must not be stored in plain text")
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
        t.Error("password hash does not match original password")
    }
}
package domain

import "time"

const (
	RoleClient     = "client"
	RoleManager    = "manager"
	RoleSuperAdmin = "super_admin"
)

type User struct {
	ID           int64
	Name         string
	Phone        string
	Email        string
	PasswordHash string
	Role         string
	TelegramID   int64
	BranchID     int64
	CreatedAt    time.Time
}

func NewClient(name, phone string) *User {
	return &User{
		Name:  name,
		Phone: phone,
		Role:  RoleClient,
	}
}

func NewManager(name, phone, email, passwordHash string, branchID int64) *User {
	return &User{
		Name:         name,
		Phone:        phone,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleManager,
		BranchID:     branchID,
	}
}

func NewSuperAdmin(name, phone, email, passwordHash string) *User {
	return &User{
		Name:         name,
		Phone:        phone,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleSuperAdmin,
	}
}

package domain

import "time"

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

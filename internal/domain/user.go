package domain

import "time"

type User struct {
	ID           int64
	Name         string
	Phone        string
	Email        string
	PasswordHash string
	Role         string
	TelegramId   int64
	BranchId     int64
	CreatedAt    time.Time
}

package domain

import "time"

type Car struct {
	ID           int64
	UserID       int64
	LicensePlate string
	Brand        string
	Model        string
	CreatedAt    time.Time
}

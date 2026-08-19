package domain

import "time"

type Car struct {
	Id int64
	UserId int64
	LicensePlate string
	Brand string
	Model string
	CreatedAt time.Time
}
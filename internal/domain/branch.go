package domain

import "time"

type Branch struct {
	Id int64
	Name string
	Address string
	Phone string
	CreatedAt time.Time
}

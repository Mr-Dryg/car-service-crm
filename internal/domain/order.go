package domain

import "time"

type Order struct {
	Id            int64
	BranchId      int64
	CarId         int64
	ServiceType   string
	Status        string
	PreferredDate string
	PreferredTime string
	Cost          float64
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

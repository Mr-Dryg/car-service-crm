package domain

import "time"

type Branch struct {
	ID        int64
	Name      string
	Address   string
	Phone     string
	CreatedAt time.Time
}

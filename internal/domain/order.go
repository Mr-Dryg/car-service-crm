package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
	StatusNew             = "new"              // Новая неподтвержденная заявка
	StatusConfirmed       = "confirmed"        // Подтвержденная заявка
	StatusCancelRequested = "cancel_requested" // Запрос на отмену
	StatusCanceled        = "canceled"         // Подтвержденная отмена
	StatusInProgress      = "in_progress"      // Ремонт
	StatusReady           = "ready"            // Машина готова, ждет выдачи
	StatusCompleted       = "completed"        // Машина выдана
)

const (
	DateLayout = "02-01-2006"
	TimeLayout = "15:04"
)

var (
	ErrInvalidStatusTransition = func(currentStatus, newStatus string) error {
		return fmt.Errorf("invalid status transition from %q to %q", currentStatus, newStatus)
	}
	ErrOrderIsCanceled  = errors.New("invalid status transition on canceled order")
	ErrOrderIsCompleted = errors.New("forbidden to change status of completed order")
)

type Order struct {
	ID              int64     `json:"id"`
	BranchID        int64     `json:"branch_id"`
	CarID           int64     `json:"car_id"`
	ServiceType     string    `json:"service_type"`
	Status          string    `json:"status"`
	PreferredDate   string    `json:"preferred_date"`
	PreferredTime   string    `json:"preferred_time"`
	Cost            float64   `json:"cost"`
	ClientConfirmed bool      `json:"client_confirmed"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (o *Order) ChangeStatus(newStatus string) error {
	var validStatuses []string

	switch o.Status {
	case StatusCancelRequested:
		validStatuses = []string{StatusConfirmed, StatusCanceled}
	case StatusCanceled:
		return ErrOrderIsCanceled
	case StatusNew:
		validStatuses = []string{StatusConfirmed, StatusCancelRequested}
	case StatusConfirmed:
		validStatuses = []string{StatusInProgress, StatusCancelRequested}
	case StatusInProgress:
		validStatuses = []string{StatusReady}
	case StatusReady:
		validStatuses = []string{StatusCompleted}
	case StatusCompleted:
		return ErrOrderIsCompleted
	}

	if slices.Contains(validStatuses, newStatus) {
		o.Status = newStatus
		return nil
	}
	return ErrInvalidStatusTransition(o.Status, newStatus)
}

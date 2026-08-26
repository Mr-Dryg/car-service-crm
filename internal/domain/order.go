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
	ID              int64
	BranchID        int64
	CarID           int64
	ServiceType     string
	Status          string
	PreferredDate   string
	PreferredTime   string
	Cost            float64
	ClientConfirmed bool
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
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

package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Mr-Dryg/car-service-crm/internal/pkg/utils"
)

const (
	StatusNew             = "new"
	StatusConfirmed       = "confirmed"
	StatusCancelRequested = "cancel_requested"
	StatusCanceled        = "canceled"
	StatusInProgress      = "in_progress"
	StatusReady           = "ready"
	StatusCompleted       = "completed"
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
	ErrOrderIsCompleted = errors.New("forbidden to update status of completed order")
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

func (o *Order) UpdateStatus(newStatus string) error {
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

type CreateOrderRequest struct {
	ClientName      string
	ClientPhone     string
	CarBrand        string
	CarModel        string
	CarLicensePlate string
	BranchID        int64
	ServiceType     string
	PreferredDate   string
	PreferredTime   string
	Notes           string
}

func (i *CreateOrderRequest) Validate() error {
	if !utils.IsValidPhone(i.ClientPhone) {
		return errors.New("invalid phone format")
	}
	if !utils.IsValidCarPlate(i.CarLicensePlate) {
		return errors.New("invalid car plate format")
	}
	return nil
}

type CreateManagerOrderRequest struct {
	CreateOrderRequest
	Status string
	Cost   float64
}

func (req *CreateManagerOrderRequest) Validate() error {
	if err := req.CreateOrderRequest.Validate(); err != nil {
		return err
	}
	return nil
}

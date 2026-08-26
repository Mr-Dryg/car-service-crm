package service

import (
	"context"
	"errors"
	"time"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

var (
	ErrInvalidDateFormat   = errors.New("invalid date format, should use DD-MM-YYYY")
	ErrDateInPast          = errors.New("forbidden to create order with a date in past")
	ErrNegativeCost        = errors.New("cost cannot be less than zero")
	ErrUpdateArchivedOrder = errors.New("forbidden to update archived order")
	ErrOrderIsNotReady     = errors.New("order must be ready to confirm by user")
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order, parsedDate time.Time) error
	GetByOrderID(ctx context.Context, orderID int64) (*domain.Order, error)
	GetByBranchID(ctx context.Context, branchID int64, includeCommon bool) ([]domain.Order, error)
	GetByCarID(ctx context.Context, carID int64) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error
	UpdateSchedule(ctx context.Context, orderID int64, prefDate time.Time, prefTime string) error
	UpdateCost(ctx context.Context, orderID int64, cost float64) error
	UpdateClientConfirmed(ctx context.Context, orderID int64, flag bool) error
	UpdateNotes(ctx context.Context, orderID int64, notes string) error
}

type OrderService struct {
	orderRepo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: repo,
	}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order) error {
	if order.BranchID < 0 {
		return errors.New("invalid branch_id")
	}
	if order.CarID <= 0 {
		return errors.New("invalid car_id")
	}
	parsedDate, err := checkTimeAvailability(order.PreferredDate)
	if err != nil {
		return err
	}

	if order.Status == "" {
		order.Status = domain.StatusNew
	}

	if order.Cost < 0 {
		return ErrNegativeCost
	}

	return s.orderRepo.Create(ctx, order, parsedDate)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	return s.orderRepo.GetByOrderID(ctx, orderID)
}

func (s *OrderService) GetBranchOrders(ctx context.Context, branchID int64, includeCommon bool) ([]domain.Order, error) {
	return s.orderRepo.GetByBranchID(ctx, branchID, includeCommon)
}

func (s *OrderService) GetCarOrders(ctx context.Context, carID int64) ([]domain.Order, error) {
	return s.orderRepo.GetByCarID(ctx, carID)
}

func (s *OrderService) ChangeStatus(ctx context.Context, orderID int64, status string) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	err = order.ChangeStatus(status)
	if err != nil {
		return err
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (s *OrderService) RescheduleOrder(ctx context.Context, orderID int64, prefDate, prefTime string) error {
	err := s.checkIfOrderArchived(ctx, orderID)
	if err != nil {
		return err
	}

	parsedDate, err := checkTimeAvailability(prefDate)
	if err != nil {
		return err
	}
	return s.orderRepo.UpdateSchedule(ctx, orderID, parsedDate, prefTime)
}

func (s *OrderService) UpdateCost(ctx context.Context, orderID int64, cost float64) error {
	err := s.checkIfOrderArchived(ctx, orderID)
	if err != nil {
		return err
	}

	if cost < 0 {
		return ErrNegativeCost
	}
	return s.orderRepo.UpdateCost(ctx, orderID, cost)
}

func (s *OrderService) ConfirmReceipt(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != domain.StatusReady {
		return ErrOrderIsNotReady
	}
	return s.orderRepo.UpdateClientConfirmed(ctx, orderID, true)
}

func (s *OrderService) UpdateNotes(ctx context.Context, orderID int64, notes string) error {
	return s.orderRepo.UpdateNotes(ctx, orderID, notes)
}

func checkTimeAvailability(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(domain.DateLayout, dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return time.Time{}, ErrDateInPast
	}
	return parsedDate, nil
}

func (s *OrderService) checkIfOrderArchived(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status == domain.StatusCompleted || order.Status == domain.StatusCanceled {
		return ErrUpdateArchivedOrder
	}
	return nil
}

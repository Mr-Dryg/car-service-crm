package service

import (
	"context"
	"errors"
	"time"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/Mr-Dryg/car-service-crm/internal/pkg/utils"
)

var (
	ErrInvalidDateFormat            = errors.New("invalid date format, should use DD-MM-YYYY")
	ErrDateInPast                   = errors.New("forbidden to create order with a date in past")
	ErrNegativeCost                 = errors.New("cost cannot be less than zero")
	ErrUpdateArchivedOrder          = errors.New("forbidden to update archived order")
	ErrOrderIsNotReady              = errors.New("order status must be ready to confirm by user")
	ErrEmptyRequiredFields          = errors.New("not all required fields are filled in")
	ErrInvalidBranchId              = errors.New("invalid branch_id")
	ErrInvalidOrderStatus           = errors.New("invalid order status")
	ErrInvalidBranchIdForConfStatus = errors.New("invalid branch_id for order with status confirmed or in_progress")
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
	orderRepo  OrderRepository
	userRepo   UserRepository
	carRepo    CarRepository
	branchRepo BranchRepository
}

func NewOrderService(or OrderRepository, ur UserRepository, cr CarRepository, br BranchRepository) *OrderService {
	return &OrderService{
		orderRepo:  or,
		userRepo:   ur,
		carRepo:    cr,
		branchRepo: br,
	}
}

type RegisterOrderInput struct {
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

type RegisterManagerOrderInput struct {
	RegisterOrderInput
	Status string
	Cost   float64
}

func (s *OrderService) CreateFromUser(ctx context.Context, input RegisterOrderInput) (*domain.Order, error) {
	return s.createOrder(ctx, input, domain.StatusNew, 0.0)
}

func (s *OrderService) CreateFromManager(ctx context.Context, input RegisterManagerOrderInput) (*domain.Order, error) {
	switch input.Status {
	case domain.StatusNew, domain.StatusConfirmed, domain.StatusInProgress, domain.StatusReady, domain.StatusCompleted, domain.StatusCanceled:
	default:
		return nil, ErrInvalidOrderStatus
	}

	if (input.Status == domain.StatusConfirmed || input.Status == domain.StatusInProgress) && input.BranchID <= 0 {
		return nil, ErrInvalidBranchIdForConfStatus
	}

	if input.Cost < 0 {
		return nil, ErrNegativeCost
	}

	return s.createOrder(ctx, input.RegisterOrderInput, input.Status, input.Cost)
}

func (s *OrderService) createOrder(ctx context.Context, input RegisterOrderInput, status string, cost float64) (*domain.Order, error) {
	parsedDate, err := checkTimeAvailability(input.BranchID, input.PreferredDate)
	if err != nil {
		return nil, err
	}

	input.ClientPhone = utils.NormalizePhone(input.ClientPhone)
	input.CarLicensePlate = utils.NormalizeCarPlate(input.CarLicensePlate)

	user, err := s.userRepo.GetByPhone(ctx, input.ClientPhone)
	if err != nil {
		user = domain.NewClient(input.ClientName, input.ClientPhone)
		if err = s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	car, err := s.carRepo.GetByLicensePlate(ctx, input.CarLicensePlate)
	if err != nil {
		car = &domain.Car{
			UserID:       user.ID,
			LicensePlate: input.CarLicensePlate,
			Brand:        input.CarBrand,
			Model:        input.CarModel,
		}
		if err = s.carRepo.Create(ctx, car); err != nil {
			return nil, err
		}
	}

	order := &domain.Order{
		BranchID:      input.BranchID,
		CarID:         car.ID,
		ServiceType:   input.ServiceType,
		Status:        status,
		PreferredDate: input.PreferredDate,
		PreferredTime: input.PreferredTime,
		Cost:          cost,
		Notes:         input.Notes,
	}
	if err = s.orderRepo.Create(ctx, order, parsedDate); err != nil {
		return nil, err
	}
	return order, nil
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
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if IsArchivedOrder(order) {
		return ErrUpdateArchivedOrder
	}

	parsedDate, err := checkTimeAvailability(order.BranchID, prefDate)
	if err != nil {
		return err
	}
	return s.orderRepo.UpdateSchedule(ctx, orderID, parsedDate, prefTime)
}

func (s *OrderService) UpdateCost(ctx context.Context, orderID int64, cost float64) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if IsArchivedOrder(order) {
		return ErrUpdateArchivedOrder
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

func checkTimeAvailability(branchID int64, dateStr string) (time.Time, error) {
	if branchID < 0 {
		return time.Time{}, ErrInvalidBranchId
	} /* else if branchID > 0 {
		TODO
	} */

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

func IsArchivedOrder(order *domain.Order) bool {
	return order.Status == domain.StatusCompleted || order.Status == domain.StatusCanceled
}

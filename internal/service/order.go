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
	ErrNegativePrice                = errors.New("price cannot be less than zero")
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
	UpdatePrice(ctx context.Context, orderID int64, price float64) error
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

func (s *OrderService) CreateFromUser(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	return s.createOrder(ctx, req, domain.StatusNew, 0.0)
}

func (s *OrderService) CreateFromManager(ctx context.Context, req domain.CreateManagerOrderRequest) (*domain.Order, error) {
	switch req.Status {
	case domain.StatusNew, domain.StatusConfirmed, domain.StatusInProgress, domain.StatusReady, domain.StatusCompleted, domain.StatusCanceled:
	default:
		return nil, ErrInvalidOrderStatus
	}

	if (req.Status == domain.StatusConfirmed || req.Status == domain.StatusInProgress) && req.BranchID <= 0 {
		return nil, ErrInvalidBranchIdForConfStatus
	}

	if req.Price < 0 {
		return nil, ErrNegativePrice
	}

	return s.createOrder(ctx, req.CreateOrderRequest, req.Status, req.Price)
}

func (s *OrderService) createOrder(ctx context.Context, req domain.CreateOrderRequest, status string, price float64) (*domain.Order, error) {
	parsedDate, err := checkTimeAvailability(req.BranchID, req.PreferredDate)
	if err != nil {
		return nil, err
	}

	req.ClientPhone = utils.NormalizePhone(req.ClientPhone)
	req.CarLicensePlate = utils.NormalizeCarPlate(req.CarLicensePlate)

	user, err := s.userRepo.GetByPhone(ctx, req.ClientPhone)
	if err != nil {
		user = domain.NewClient(req.ClientName, req.ClientPhone)
		if err = s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	car, err := s.carRepo.GetByLicensePlate(ctx, req.CarLicensePlate)
	if err != nil {
		car = &domain.Car{
			UserID:       user.ID,
			LicensePlate: req.CarLicensePlate,
			Brand:        req.CarBrand,
			Model:        req.CarModel,
		}
		if err = s.carRepo.Create(ctx, car); err != nil {
			return nil, err
		}
	}

	order := &domain.Order{
		BranchID:      req.BranchID,
		CarID:         car.ID,
		ServiceType:   req.ServiceType,
		Status:        status,
		PreferredDate: req.PreferredDate,
		PreferredTime: req.PreferredTime,
		Price:         price,
		Notes:         req.Notes,
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

func (s *OrderService) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	err = order.UpdateStatus(status)
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

	if isArchivedOrder(order) {
		return ErrUpdateArchivedOrder
	}

	parsedDate, err := checkTimeAvailability(order.BranchID, prefDate)
	if err != nil {
		return err
	}
	return s.orderRepo.UpdateSchedule(ctx, orderID, parsedDate, prefTime)
}

func (s *OrderService) UpdatePrice(ctx context.Context, orderID int64, price float64) error {
	order, err := s.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if isArchivedOrder(order) {
		return ErrUpdateArchivedOrder
	}

	if price < 0 {
		return ErrNegativePrice
	}
	return s.orderRepo.UpdatePrice(ctx, orderID, price)
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

func isArchivedOrder(order *domain.Order) bool {
	return order.Status == domain.StatusCompleted || order.Status == domain.StatusCanceled
}

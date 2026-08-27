package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/Mr-Dryg/car-service-crm/internal/pkg/utils"
	"github.com/Mr-Dryg/car-service-crm/internal/service"
)

type OrderService interface {
	CreateFromUser(ctx context.Context, input service.RegisterOrderInput) (*domain.Order, error)
	CreateFromManager(ctx context.Context, input service.RegisterManagerOrderInput) (*domain.Order, error)
}

type OrderHandler struct {
	service OrderService
}

type CreateUserOrderRequest struct {
	service.RegisterOrderInput
}

func (req *CreateUserOrderRequest) Validate() error {
	if !utils.IsValidPhone(req.ClientPhone) {
		return errors.New("invalid phone format")
	}
	if !utils.IsValidCarPlate(req.CarLicensePlate) {
		return errors.New("invalid car plate format")
	}
	return nil
}

type CreateManagerOrderRequest struct {
	service.RegisterManagerOrderInput
}

func (req *CreateManagerOrderRequest) Validate() error {
	userReq := CreateUserOrderRequest{RegisterOrderInput: req.RegisterOrderInput}
	if err := userReq.Validate(); err != nil {
		return err
	}
	return nil
}

func NewOrderHandler(srv OrderService) *OrderHandler {
	return &OrderHandler{
		service: srv,
	}
}

func (h *OrderHandler) CreateFromUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateUserOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateFromUser(r.Context(), req.RegisterOrderInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "reques coding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) CreateFromManager(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateManagerOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateFromManager(r.Context(), req.RegisterManagerOrderInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "reques coding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

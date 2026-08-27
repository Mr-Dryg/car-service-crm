package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

type OrderService interface {
	CreateFromUser(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error)
	CreateFromManager(ctx context.Context, req domain.CreateManagerOrderRequest) (*domain.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(srv OrderService) *OrderHandler {
	return &OrderHandler{
		service: srv,
	}
}

func (h *OrderHandler) CreateFromUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req domain.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateFromUser(r.Context(), req)
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

	var req domain.CreateManagerOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateFromManager(r.Context(), req)
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

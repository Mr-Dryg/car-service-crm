package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

type OrderService interface {
	Create(ctx context.Context, order *domain.Order) error
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(srv OrderService) *OrderHandler {
	return &OrderHandler{
		service: srv,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var order domain.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &order); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "reques coding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

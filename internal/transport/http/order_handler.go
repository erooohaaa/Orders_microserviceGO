package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"Orders/internal/usecase"
)

type OrderHandler struct {
	uc *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.createOrder)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
	mux.HandleFunc("PATCH /orders/{id}/cancel", h.cancelOrder)
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string `json:"customer_id"`
		ItemName   string `json:"item_name"`
		Amount     int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.CreateOrder(r.Context(), req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		if errors.Is(err, usecase.ErrGatewayUnavailable) {
			http.Error(w, "payment service unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.uc.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.CancelOrder(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

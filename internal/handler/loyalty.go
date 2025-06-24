package handler

import (
	"encoding/json"
	"gophermart/internal/errors"
	"gophermart/internal/middleware"
	"gophermart/internal/service"
	"net/http"
)

type loyaltyHandler struct {
	service service.Service
}

func NewLoyaltyHandler(service service.Service) *loyaltyHandler {
	return &loyaltyHandler{
		service: service,
	}
}

func getUserIDFromContext(r *http.Request) int64 {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		return 0
	}
	return userID
}

func (h *loyaltyHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Order string `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if req.Order == "" {
		http.Error(w, "order is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.UploadOrder(r.Context(), userID, req.Order)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalidOrderNumberFormat):
			http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		case errors.Is(err, errors.ErrOrderAlreadyUploadedByAnotherUser):
			http.Error(w, "order already uploaded by another user", http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(order)
}

func (h *loyaltyHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *loyaltyHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

func (h *loyaltyHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if req.Order == "" || req.Sum <= 0 {
		http.Error(w, "order and sum are required", http.StatusBadRequest)
		return
	}

	withdrawal, err := h.service.Withdraw(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalidOrderNumberFormat):
			http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		case errors.Is(err, errors.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawal)
}

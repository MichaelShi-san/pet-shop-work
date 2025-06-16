package handlers

import (
	"encoding/json"
	"go-pet-shop/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Orders interface {
	CreateOrder(order models.Order) (int, error)
	GetOrderByID(id int) (models.Order, error)
	AddOrderItem(item models.OrderItem) error
	GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error)
	GetOrdersByUserEmail(email string) ([]models.Order, error)
}

type OrdersHandler struct {
	log *slog.Logger
	Storage Orders
}

func NewOrdersHandler(log *slog.Logger, storage Orders) *OrdersHandler {
	return &OrdersHandler{
		log:     log,
		Storage: storage,
	}
}

func (h *OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Storage.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *OrdersHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item models.OrderItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		h.log.Error("failed to decode request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.OrderID = orderID

	if err := h.Storage.AddOrderItem(item); err != nil {
		h.log.Error("failed to add order item", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Order item added successfully"})
}
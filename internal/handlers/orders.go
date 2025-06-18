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
	PlaceOrder(userEmail string, items []models.OrderItem) (int, error)
	GetUserOrderHistory(email string) ([]models.OrderDetail, error)
}

type OrdersHandler struct {
	log *slog.Logger
	Storage Orders
	
}

type orderRequest struct {
    UserEmail string             `json:"user_email"`
    Items     []models.OrderItem `json:"items"`
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

func (h *OrdersHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
    var req orderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    h.log.Info("PlaceOrder request", slog.Any("user_email", req.UserEmail), slog.Any("items", req.Items))

    orderID, err := h.Storage.PlaceOrder(req.UserEmail, req.Items)
    if err != nil {
        h.log.Error("failed to place order", slog.Any("error", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
}

func (h *OrdersHandler) GetUserOrderHistory(w http.ResponseWriter, r *http.Request) {
    email := chi.URLParam(r, "email")
    if email == "" {
        http.Error(w, "email is required", http.StatusBadRequest)
        return
    }

    history, err := h.Storage.GetUserOrderHistory(email)
    if err != nil {
        h.log.Error("failed to get user order history", slog.Any("error", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(history)
}




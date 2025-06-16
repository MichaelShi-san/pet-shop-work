package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-pet-shop/models"
)

type PostgresStorage struct {
	db *sql.DB
}

type Storage interface {
	GetAllProducts() ([]models.Product, error)
	CreateProduct(product models.Product) error
	UpdateProduct(product models.Product) error
	DeleteProduct(id int) error

	CreateOrder(order models.Order) (int, error)
	AddOrderItem(item models.OrderItem) error
	GetOrderByID(id int) (models.Order, error)
	GetOrdersByUserEmail(email string) ([]models.Order, error)
	GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) 
	}

	func (s *PostgresStorage) PlaceOrder(userEmail string, items []models.OrderItem) (orderID int, err error) {
    tx, err := s.db.Begin()
    if err != nil {
        return 0, err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // Получаем user_id по email
    var userID int
    err = tx.QueryRow("SELECT id FROM users WHERE email = $1", userEmail).Scan(&userID)
    if err != nil {
        return 0, fmt.Errorf("user not found: %w", err)
    }

    // Проверяем и уменьшаем stock для каждого товара
    for _, item := range items {
        res, err := tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1", item.Quantity, item.ProductID)
        if err != nil {
            return 0, fmt.Errorf("failed to update stock: %w", err)
        }
        rowsAffected, err := res.RowsAffected()
        if err != nil || rowsAffected == 0 {
            return 0, fmt.Errorf("not enough stock for product %d", item.ProductID)
        }
    }

    // Создаем заказ
    totalPrice := 0.0
	
    // Можно посчитать сумму или сделать отдельный запрос
    for _, item := range items {
        var price float64
        err = tx.QueryRow("SELECT price FROM products WHERE id = $1", item.ProductID).Scan(&price)
        if err != nil {
            return 0, err
        }
        totalPrice += price * float64(item.Quantity)
    }

    err = tx.QueryRow(
        "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id",
        userID, totalPrice).Scan(&orderID)
    if err != nil {
        return 0, err
    }

    // Добавляем позиции заказа
    for _, item := range items {
        _, err = tx.Exec(
            "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)",
            orderID, item.ProductID, item.Quantity)
        if err != nil {
            return 0, err
        }
    }

    // Создаем транзакцию оплаты
    _, err = tx.Exec(
        "INSERT INTO transactions (order_id, amount, status, created_at) VALUES ($1, $2, $3, NOW())",
        orderID, totalPrice, "completed") // или статус "pending"
    if err != nil {
        return 0, err
    }

    err = tx.Commit()
    if err != nil {
        return 0, err
    }

    return orderID, nil
}


var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
)

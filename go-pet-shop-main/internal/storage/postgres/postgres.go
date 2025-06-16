package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-pet-shop/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// PlaceOrder implements handlers.Orders.
func (s *Storage) PlaceOrder(userEmail string, items []models.OrderItem) (int, error) {
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Найти user_id по email
	var userID int
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, userEmail).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	// Создать заказ
	var orderID int
	err = tx.QueryRow(ctx, `INSERT INTO orders (user_id, total_price) VALUES ($1, 0) RETURNING id`, userID).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	// 3Добавить товары
	var total float64
	for _, item := range items {
		var price float64
		err = tx.QueryRow(ctx, `SELECT price FROM products WHERE id = $1`, item.ProductID).Scan(&price)
		if err != nil {
			return 0, fmt.Errorf("product not found: %w", err)
		}

		_, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			orderID, item.ProductID, item.Quantity)
		if err != nil {
			return 0, fmt.Errorf("failed to insert order item: %w", err)
		}

		total += price * float64(item.Quantity)
	}

	// Обновить общую сумму
	_, err = tx.Exec(ctx, `UPDATE orders SET total_price = $1 WHERE id = $2`, total, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to update total price: %w", err)
	}

	// Записать транзакцию
	_, err = tx.Exec(ctx, `INSERT INTO transactions (order_id, amount, status) VALUES ($1, $2, $3)`,
		orderID, total, "pending")
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit: %w", err)
	}

	return orderID, nil
}

func New(databaseUrl string) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	s.db.Close()
	return nil
}

func (s *Storage) CreateOrder(order models.Order) (int, error) {
	query := `
		INSERT INTO orders (user_id, total_price)
		VALUES ($1, 0)
		RETURNING id;
	`

	var id int
	err := s.db.QueryRow(context.Background(), query, order.CustomerID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	return id, nil
}

func (s *Storage) GetOrderByID(id int) (models.Order, error) {
	query := `
		SELECT id, user_id, created_at
		FROM orders
		WHERE id = $1;
	`

	var order models.Order
	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, nil
		}
		return order, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (s *Storage) GetOrdersByUserEmail(email string) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.created_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE u.email = $1;
	`

	rows, err := s.db.Query(context.Background(), query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity
		FROM order_items
		WHERE order_id = $1;
	`

	rows, err := s.db.Query(context.Background(), query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *Storage) AddOrderItem(item models.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity)
		VALUES ($1, $2, $3);
	`

	_, err := s.db.Exec(context.Background(), query, item.OrderID, item.ProductID, item.Quantity)
	if err != nil {
		return fmt.Errorf("add order item: %w", err)
	}
	return nil
}

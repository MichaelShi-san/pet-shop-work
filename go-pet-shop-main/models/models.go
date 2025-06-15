package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type Customer struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	CustomerID int
	CreatedAt  time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
}

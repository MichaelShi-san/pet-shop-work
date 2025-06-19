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
    ID        int `json:"id"`
    OrderID   int `json:"order_id"`
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}

type OrderDetail struct {
	OrderID int `json:"order_id"`
	ProductName string `json:"product_name"`
	ProductID int `json:"product_id"`
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	TransactionStatus string `json:"transaction_status"`
}

type PopularProduct struct {
	ID int `json:"id"`
	Name string `json:"name"`
	TotalSold int `json:"total_sold"`
}

type User struct {
	ID    int
	Name  string
	Email string
}

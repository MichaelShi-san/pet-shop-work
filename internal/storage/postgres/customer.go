package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/models"
)

func (s *Storage) CreateCustomer(c models.Customer) error {
	const fn = "storage.postgres.customer.CreateCustomer"

	_, err := s.db.Exec(context.Background(),
		`INSERT INTO customers (name, email) VALUES ($1, $2)`,
		c.Name, c.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetCustomerByEmail(email string) (models.Customer, error) {
	const fn = "storage.postgres.customer.GetCustomerByEmail"

	var c models.Customer
	err := s.db.QueryRow(context.Background(),
		`SELECT id, name, email FROM customers WHERE email = $1`, email).
		Scan(&c.ID, &c.Name, &c.Email)
	if err != nil {
		return c, fmt.Errorf("%s: %w", fn, err)
	}

	return c, nil
}

func (s *Storage) GetAllCustomers() ([]models.Customer, error) {
	const fn = "storage.postgres.customer.GetAllCustomers"

	rows, err := s.db.Query(context.Background(), `SELECT id, name, email FROM customers`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		customers = append(customers, c)
	}

	return customers, nil
}


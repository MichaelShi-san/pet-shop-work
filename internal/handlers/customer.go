package handlers

import (
	"go-pet-shop/internal/storage"
	"go-pet-shop/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Customers interface {
	CreateCustomer(models.Customer) error
	GetCustomerByEmail(string) (models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
}

// type Users interface {
// 	CreateUser(user models.User) error
// 	GetUserByEmail(email string) (models.User, error)
// 	GetAllUsers() ([]models.User, error)
// }


func CreateCustomer(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.customers.CreateCustomer"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var customer models.Customer
		if err := render.DecodeJSON(r.Body, &customer); err != nil {
			log.Error("decode error", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := customers.CreateCustomer(customer); err != nil {
			log.Error("create error", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, map[string]string{"status": "Customer created"})
	}
}

func GetCustomerByEmail(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.customers.GetCustomerByEmail"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := chi.URLParam(r, "email")
		if email == "" {
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}

		customer, err := customers.GetCustomerByEmail(email)
		if err != nil {
			log.Error("get error", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, customer)
	}
}

func GetAllCustomers(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.customers.GetAllCustomers"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		all, err := customers.GetAllCustomers()
		if err != nil {
			log.Error("get all error", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, all)
	}
}

func GetUserByEmail(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		if email == "" {
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}

		user, err := storage.GetUserByEmail(email)
		if err != nil {
			log.Error("failed to get user by email", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, user)
	}
}

func CreateUser(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.CreateUser"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user models.User
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := storage.CreateUser(user); err != nil {
			log.Error("failed to create user", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, map[string]string{"status": "user created"})
	}
}

func GetAllUsers(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.GetAllUsers"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		usersList, err := storage.GetAllUsers()
		if err != nil {
			log.Error("failed to get all users", slog.Any("err", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, usersList)
	}
}


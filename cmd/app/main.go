package main

import (
	"go-pet-shop/internal/config"
	"go-pet-shop/internal/handlers"
	"go-pet-shop/internal/lib/logger"
	"go-pet-shop/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting the project...", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.CustomLogger(log))

	router.Get("/health", handlers.StatusHandler)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", handlers.GetAllProducts(log, storage))
		r.Post("/", handlers.CreateProduct(log, storage))
		r.Get("/{id}", handlers.GetProductByID(log, storage))
		r.Put("/{id}", handlers.UpdateProduct(log, storage))
		r.Delete("/{id}", handlers.DeleteProduct(log, storage))

		r.Get("/popular", handlers.GetPopularProducts(log, storage))


	})

	ordersHandler := handlers.NewOrdersHandler(log, storage)
	router.Route("/orders", func(r chi.Router) {
		r.Post("/", ordersHandler.CreateOrder)
		r.Post("/{id}/items", ordersHandler.AddOrderItem)
		r.Post("/place", ordersHandler.PlaceOrder)
	})

	router.Route("/users", func(r chi.Router) {
		r.Get("/{email}/history", ordersHandler.GetUserOrderHistory)
	})

	handler := logger.LoggingMiddleware(log, router)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Starting server on", slog.String("address", cfg.HTTPServer.Address))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server error: ", slog.String("err", err.Error()))
	}
}

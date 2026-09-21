package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cassiokiyoshi/cheap-eats/internal/database"
	"github.com/cassiokiyoshi/cheap-eats/internal/handlers"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	databasePool, err := database.Open(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer databasePool.Close()

	log.Println("Connected to PostgreSQL")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",
			"http://127.0.0.1:8081",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	dishRepository :=
		repository.NewPostgresDishRepository(databasePool)

	restaurantRepository :=
		repository.NewPostgresRestaurantRepository(databasePool)

	dishHandler := handlers.NewDishHandler(
		dishRepository,
		restaurantRepository,
	)

	dishSearchHandler := handlers.NewDishSearchHandler(dishRepository)

	restaurantHandler := handlers.NewRestaurantHandler(
		restaurantRepository,
	)

	router.Get("/api/health", handlers.Health)

	router.Get("/api/dishes", dishHandler.List)
	router.Get("/api/dishes/nearby", dishSearchHandler.Nearby)
	router.Get("/api/dishes/{id}", dishHandler.Get)
	router.Post("/api/dishes", dishHandler.Create)
	router.Patch("/api/dishes/{id}/price", dishHandler.UpdatePrice)
	router.Get(
		"/api/dishes/{id}/price-history",
		dishHandler.PriceHistory,
	)

	router.Get("/api/restaurants", restaurantHandler.List)
	router.Get("/api/restaurants/nearby", restaurantHandler.Nearby)
	router.Get("/api/restaurants/{id}", restaurantHandler.Get)
	router.Post("/api/restaurants", restaurantHandler.Create)
	router.Get(
		"/api/restaurants/{id}/dishes",
		dishHandler.ListByRestaurant,
	)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stopContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting Cheap Eats API at http://localhost%s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server: %w", err)

	case <-stopContext.Done():
		log.Println("Shutting down HTTP server...")
	}

	// Restore normal signal handling so a second Ctrl+C forces exit.
	stop()

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		// The grace period expired, or shutdown failed.
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("Force-close HTTP server: %v", closeErr)
		}
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	log.Println("HTTP server stopped")
	return nil
}

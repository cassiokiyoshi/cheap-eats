package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cassiokiyoshi/cheap-eats/internal/database"
	"github.com/cassiokiyoshi/cheap-eats/internal/handlers"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/cassiokiyoshi/cheap-eats/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	databasePool, err := database.Open(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer databasePool.Close()

	log.Println("Connected to PostgreSQL")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	dishRepository := repository.NewDishRepository()
	restaurantRepository := repository.NewRestaurantRepository()

	dishHandler := handlers.NewDishHandler(
		dishRepository,
		restaurantRepository,
	)

	dishService := service.NewDishService(
		dishRepository,
		restaurantRepository,
	)
	dishSearchHandler := handlers.NewDishSearchHandler(dishService)

	restaurantHandler := handlers.NewRestaurantHandler(
		restaurantRepository,
	)

	router.Get("/api/health", handlers.Health)

	router.Get("/api/dishes", dishHandler.List)
	router.Get("/api/dishes/nearby", dishSearchHandler.Nearby)
	router.Get("/api/dishes/{id}", dishHandler.Get)
	router.Post("/api/dishes", dishHandler.Create)

	router.Get("/api/restaurants", restaurantHandler.List)
	router.Get("/api/restaurants/nearby", restaurantHandler.Nearby)
	router.Get("/api/restaurants/{id}", restaurantHandler.Get)

	address := ":8080"
	fmt.Printf("Cheap Eats API running at http://localhost%s\n", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}

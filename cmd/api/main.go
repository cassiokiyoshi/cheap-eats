package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cassiokiyoshi/cheap-eats/internal/handlers"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	dishRepository := repository.NewDishRepository()
	dishHandler := handlers.NewDishHandler(dishRepository)

	router.Get("/api/health", handlers.Health)
	router.Get("/api/dishes", dishHandler.List)
	router.Get("/api/dishes/{id}", dishHandler.Get)

	address := ":8080"
	fmt.Printf("Cheap Eats API running at http://localhost%s\n", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
)

type DishHandler struct {
	dishRepository       repository.DishStore
	restaurantRepository repository.RestaurantStore
}

func NewDishHandler(
	dishRepository repository.DishStore,
	restaurantRepository repository.RestaurantStore,
) *DishHandler {
	return &DishHandler{
		dishRepository:       dishRepository,
		restaurantRepository: restaurantRepository,
	}
}

func (h *DishHandler) List(w http.ResponseWriter, r *http.Request) {
	var (
		dishes []models.Dish
		err    error
	)

	if value := r.URL.Query().Get("max_price"); value != "" {
		maxPrice, parseErr := strconv.Atoi(value)
		if parseErr != nil || maxPrice <= 0 {
			http.Error(
				w,
				"max_price must be a positive integer",
				http.StatusBadRequest,
			)
			return
		}

		dishes, err = h.dishRepository.ListByMaxPrice(
			r.Context(),
			maxPrice,
		)
	} else {
		dishes, err = h.dishRepository.List(r.Context())
	}

	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dishes)
}

func (h *DishHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid dish ID", http.StatusBadRequest)
		return
	}

	dish, found, err := h.dishRepository.FindByID(
		r.Context(),
		id,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if !found {
		http.Error(w, "dish not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dish)
}

type createDishInput struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Currency     string `json:"currency"`
	RestaurantID int64  `json:"restaurant_id"`
}

func (h *DishHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input createDishInput

	if err := readJSON(w, r, &input); err != nil {
		var sizeError *http.MaxBytesError

		if errors.As(err, &sizeError) {
			http.Error(
				w,
				"request body must not exceed 64 KiB",
				http.StatusRequestEntityTooLarge,
			)
			return
		}

		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Currency = strings.ToUpper(strings.TrimSpace(input.Currency))

	if input.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if input.RestaurantID < 1 {
		http.Error(w, "restaurant_id is required", http.StatusBadRequest)
		return
	}

	_, found, err := h.restaurantRepository.FindByID(
		r.Context(),
		input.RestaurantID,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if !found {
		http.Error(w, "restaurant not found", http.StatusBadRequest)
		return
	}

	if input.Price <= 0 {
		http.Error(w, "price must be greater than zero", http.StatusBadRequest)
		return
	}

	if input.Currency == "" {
		input.Currency = "JPY"
	}

	dish := models.Dish{
		Name:         input.Name,
		Price:        input.Price,
		Currency:     input.Currency,
		RestaurantID: input.RestaurantID,
	}

	createdDish, err := h.dishRepository.Create(
		r.Context(),
		dish,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDish)
}

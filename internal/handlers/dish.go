package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
)

type DishHandler struct {
	dishRepository       *repository.DishRepository
	restaurantRepository *repository.RestaurantRepository
}

func NewDishHandler(
	dishRepository *repository.DishRepository,
	restaurantRepository *repository.RestaurantRepository,
) *DishHandler {
	return &DishHandler{
		dishRepository:       dishRepository,
		restaurantRepository: restaurantRepository,
	}
}

func (h *DishHandler) List(w http.ResponseWriter, r *http.Request) {
	dishes := h.dishRepository.List()

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

	dish, found := h.dishRepository.FindByID(id)
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
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

	if _, found := h.restaurantRepository.FindByID(input.RestaurantID); !found {
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

	createdDish := h.dishRepository.Create(dish)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDish)
}

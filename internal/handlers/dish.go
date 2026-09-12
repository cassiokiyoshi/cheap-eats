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
	repository *repository.DishRepository
}

func NewDishHandler(repository *repository.DishRepository) *DishHandler {
	return &DishHandler{
		repository: repository,
	}
}

func (h *DishHandler) List(w http.ResponseWriter, r *http.Request) {
	dishes := h.repository.List()

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

	dish, found := h.repository.FindByID(id)
	if !found {
		http.Error(w, "dish not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dish)
}

type createDishInput struct {
	Name           string `json:"name"`
	Price          int    `json:"price"`
	Currency       string `json:"currency"`
	RestaurantName string `json:"restaurant_name"`
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
	input.RestaurantName = strings.TrimSpace(input.RestaurantName)

	if input.Name == "" || input.RestaurantName == "" {
		http.Error(w, "name and restaurant_name are required", http.StatusBadRequest)
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
		Name:           input.Name,
		Price:          input.Price,
		Currency:       input.Currency,
		RestaurantName: input.RestaurantName,
	}

	createdDish := h.repository.Create(dish)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDish)
}

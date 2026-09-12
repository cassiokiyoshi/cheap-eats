package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
)

type RestaurantHandler struct {
	repository *repository.RestaurantRepository
}

func NewRestaurantHandler(
	repository *repository.RestaurantRepository,
) *RestaurantHandler {
	return &RestaurantHandler{
		repository: repository,
	}
}

func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurants := h.repository.List()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurants)
}

func (h *RestaurantHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid restaurant ID", http.StatusBadRequest)
		return
	}

	restaurant, found := h.repository.FindByID(id)
	if !found {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurant)
}

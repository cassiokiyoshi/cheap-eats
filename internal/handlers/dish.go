package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

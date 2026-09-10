package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
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

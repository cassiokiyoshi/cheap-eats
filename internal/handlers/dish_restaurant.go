package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *DishHandler) ListByRestaurant(
	w http.ResponseWriter,
	r *http.Request,
) {
	restaurantID, err := strconv.ParseInt(
		chi.URLParam(r, "id"),
		10,
		64,
	)
	if err != nil || restaurantID < 1 {
		http.Error(w, "invalid restaurant ID", http.StatusBadRequest)
		return
	}

	_, found, err := h.restaurantRepository.FindByID(
		r.Context(),
		restaurantID,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if !found {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	dishes, err := h.dishRepository.ListByRestaurantID(
		r.Context(),
		restaurantID,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(dishes); err != nil {
		log.Printf("encode restaurant dishes response: %v", err)
	}
}

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type updateDishPriceInput struct {
	Price int `json:"price"`
}

func (h *DishHandler) UpdatePrice(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid dish ID", http.StatusBadRequest)
		return
	}

	var input updateDishPriceInput

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

	// PostgreSQL INTEGER is a signed 32-bit integer.
	if err := validatePrice(input.Price); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dish, found, err := h.dishRepository.UpdatePrice(
		r.Context(),
		id,
		input.Price,
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

	if err := json.NewEncoder(w).Encode(dish); err != nil {
		log.Printf("encode updated dish response: %v", err)
	}
}

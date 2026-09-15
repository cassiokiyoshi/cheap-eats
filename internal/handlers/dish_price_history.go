package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *DishHandler) PriceHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid dish ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	limit := 20
	if query.Has("limit") {
		value, err := strconv.Atoi(query.Get("limit"))
		if err != nil || value < 1 || value > 100 {
			http.Error(
				w,
				"limit must be between 1 and 100",
				http.StatusBadRequest,
			)
			return
		}
		limit = value
	}

	offset := 0
	if query.Has("offset") {
		value, err := strconv.Atoi(query.Get("offset"))
		if err != nil || value < 0 {
			http.Error(
				w,
				"offset must be a nonnegative integer",
				http.StatusBadRequest,
			)
			return
		}
		offset = value
	}

	_, found, err := h.dishRepository.FindByID(r.Context(), id)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if !found {
		http.Error(w, "dish not found", http.StatusNotFound)
		return
	}

	history, err := h.dishRepository.ListPriceHistory(
		r.Context(),
		id,
		limit,
		offset,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(history); err != nil {
		log.Printf("encode dish price history response: %v", err)
	}
}

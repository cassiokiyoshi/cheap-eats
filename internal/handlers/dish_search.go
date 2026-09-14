package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type DishSearcher interface {
	SearchNearby(
		ctx context.Context,
		latitude float64,
		longitude float64,
		radiusMeters float64,
		maxPrice int,
		sortBy string,
		limit int,
		offset int,
	) ([]models.DishSearchResult, error)
}

type DishSearchHandler struct {
	service DishSearcher
}

func NewDishSearchHandler(
	service DishSearcher,
) *DishSearchHandler {
	return &DishSearchHandler{
		service: service,
	}
}

func (h *DishSearchHandler) Nearby(
	w http.ResponseWriter,
	r *http.Request,
) {
	latitude, err := strconv.ParseFloat(
		r.URL.Query().Get("lat"),
		64,
	)
	if err != nil || math.IsNaN(latitude) || math.IsInf(latitude, 0) ||
		latitude < -90 || latitude > 90 {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(
		r.URL.Query().Get("lng"),
		64,
	)
	if err != nil || math.IsNaN(longitude) || math.IsInf(longitude, 0) ||
		longitude < -180 || longitude > 180 {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	radius := 1000.0
	if value := r.URL.Query().Get("radius"); value != "" {
		radius, err = strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(radius) || math.IsInf(radius, 0) ||
			radius < 1 || radius > 5000 {
			http.Error(
				w,
				"radius must be between 1 and 5000 meters",
				http.StatusBadRequest,
			)
			return
		}
	}

	maxPrice, err := strconv.Atoi(
		r.URL.Query().Get("max_price"),
	)
	if err != nil || maxPrice <= 0 {
		http.Error(
			w,
			"max_price must be a positive integer",
			http.StatusBadRequest,
		)
		return
	}

	sortBy := r.URL.Query().Get("sort")

	switch sortBy {
	case "", "distance":
		sortBy = "distance"
	case "price":
		// Valid option.
	default:
		http.Error(
			w,
			"sort must be distance or price",
			http.StatusBadRequest,
		)
		return
	}

	limit := 20
	if values, exists := r.URL.Query()["limit"]; exists {
		value, parseErr := strconv.Atoi(values[0])
		if parseErr != nil || value < 1 || value > 100 {
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
	if values, exists := r.URL.Query()["offset"]; exists {
		value, parseErr := strconv.Atoi(values[0])
		if parseErr != nil || value < 0 {
			http.Error(
				w,
				"offset must be a nonnegative integer",
				http.StatusBadRequest,
			)
			return
		}
		offset = value
	}

	results, err := h.service.SearchNearby(
		r.Context(),
		latitude,
		longitude,
		radius,
		maxPrice,
		sortBy,
		limit,
		offset,
	)

	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

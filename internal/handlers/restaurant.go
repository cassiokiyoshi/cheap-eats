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

func (h *RestaurantHandler) Nearby(
	w http.ResponseWriter,
	r *http.Request,
) {
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	radius := 500.0

	if value := r.URL.Query().Get("radius"); value != "" {
		radius, err = strconv.ParseFloat(value, 64)
		if err != nil || radius <= 0 || radius > 5000 {
			http.Error(
				w,
				"radius must be between 1 and 5000 meters",
				http.StatusBadRequest,
			)
			return
		}
	}

	restaurants := h.repository.Nearby(
		latitude,
		longitude,
		radius,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurants)
}

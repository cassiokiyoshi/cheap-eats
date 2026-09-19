package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type createRestaurantInput struct {
	Name      string   `json:"name"`
	NameJA    *string  `json:"name_ja"`
	NameEN    *string  `json:"name_en"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (h *RestaurantHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input createRestaurantInput

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

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	if input.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	nameJA, err := normalizeOptionalName(input.NameJA, "name_ja")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nameEN, err := normalizeOptionalName(input.NameEN, "name_en")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}

	if input.Latitude == nil {
		http.Error(w, "latitude is required", http.StatusBadRequest)
		return
	}

	if input.Longitude == nil {
		http.Error(w, "longitude is required", http.StatusBadRequest)
		return
	}

	latitude := *input.Latitude
	longitude := *input.Longitude

	if math.IsNaN(latitude) || math.IsInf(latitude, 0) ||
		latitude < -90 || latitude > 90 {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}

	if math.IsNaN(longitude) || math.IsInf(longitude, 0) ||
		longitude < -180 || longitude > 180 {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	restaurant := models.Restaurant{
		Name:      input.Name,
		NameJA:    nameJA,
		NameEN:    nameEN,
		Address:   input.Address,
		Latitude:  latitude,
		Longitude: longitude,
	}

	createdRestaurant, err := h.repository.Create(
		r.Context(),
		restaurant,
	)
	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(
		"Location",
		"/api/restaurants/"+strconv.FormatInt(createdRestaurant.ID, 10),
	)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdRestaurant); err != nil {
		// The response has already started, so don't write another status.
		return
	}
}

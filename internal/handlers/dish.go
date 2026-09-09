package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

func ListDishes(w http.ResponseWriter, r *http.Request) {
	dishes := []models.Dish{
		{
			ID:             1,
			Name:           "Shoyu Ramen",
			Price:          850,
			Currency:       "JPY",
			RestaurantName: "Tokyo Ramen",
		},
		{
			ID:             2,
			Name:           "Gyudon",
			Price:          650,
			Currency:       "JPY",
			RestaurantName: "Cheap Bowl",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dishes)
}

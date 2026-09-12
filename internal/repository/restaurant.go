package repository

import (
	"sync"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type RestaurantRepository struct {
	mu          sync.RWMutex
	restaurants []models.Restaurant
}

func NewRestaurantRepository() *RestaurantRepository {
	return &RestaurantRepository{
		restaurants: []models.Restaurant{
			{
				ID:        1,
				Name:      "Tokyo Ramen",
				Address:   "Marunouchi, Tokyo",
				Latitude:  35.6812,
				Longitude: 139.7671,
			},
			{
				ID:        2,
				Name:      "Cheap Bowl",
				Address:   "Kanda, Tokyo",
				Latitude:  35.6917,
				Longitude: 139.7709,
			},
		},
	}
}

func (r *RestaurantRepository) List() []models.Restaurant {
	r.mu.RLock()
	defer r.mu.RUnlock()

	restaurants := make([]models.Restaurant, len(r.restaurants))
	copy(restaurants, r.restaurants)

	return restaurants
}

func (r *RestaurantRepository) FindByID(id int64) (models.Restaurant, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, restaurant := range r.restaurants {
		if restaurant.ID == id {
			return restaurant, true
		}
	}

	return models.Restaurant{}, false
}

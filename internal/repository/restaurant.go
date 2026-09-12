package repository

import (
	"math"
	"sort"
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

func (r *RestaurantRepository) Nearby(
	latitude float64,
	longitude float64,
	radiusMeters float64,
) []models.RestaurantSuggestion {
	r.mu.RLock()
	defer r.mu.RUnlock()

	suggestions := make([]models.RestaurantSuggestion, 0)

	for _, restaurant := range r.restaurants {
		distance := distanceMeters(
			latitude,
			longitude,
			restaurant.Latitude,
			restaurant.Longitude,
		)

		if distance <= radiusMeters {
			suggestions = append(
				suggestions,
				models.RestaurantSuggestion{
					Restaurant:     restaurant,
					DistanceMeters: int(math.Round(distance)),
				},
			)
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].DistanceMeters <
			suggestions[j].DistanceMeters
	})

	return suggestions
}

func distanceMeters(
	latitude1 float64,
	longitude1 float64,
	latitude2 float64,
	longitude2 float64,
) float64 {
	const earthRadiusMeters = 6_371_000

	lat1 := latitude1 * math.Pi / 180
	lat2 := latitude2 * math.Pi / 180
	latDifference := (latitude2 - latitude1) * math.Pi / 180
	lngDifference := (longitude2 - longitude1) * math.Pi / 180

	a := math.Sin(latDifference/2)*math.Sin(latDifference/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(lngDifference/2)*math.Sin(lngDifference/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

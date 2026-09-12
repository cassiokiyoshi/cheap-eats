package service

import (
	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
)

type DishService struct {
	dishRepository       *repository.DishRepository
	restaurantRepository *repository.RestaurantRepository
}

func NewDishService(
	dishRepository *repository.DishRepository,
	restaurantRepository *repository.RestaurantRepository,
) *DishService {
	return &DishService{
		dishRepository:       dishRepository,
		restaurantRepository: restaurantRepository,
	}
}

func (s *DishService) SearchNearby(
	latitude float64,
	longitude float64,
	radiusMeters float64,
	maxPrice int,
) []models.DishSearchResult {
	restaurants := s.restaurantRepository.Nearby(
		latitude,
		longitude,
		radiusMeters,
	)
	dishes := s.dishRepository.ListByMaxPrice(maxPrice)

	results := make([]models.DishSearchResult, 0)

	for _, restaurant := range restaurants {
		for _, dish := range dishes {
			if dish.RestaurantID == restaurant.ID {
				results = append(results, models.DishSearchResult{
					Dish:           dish,
					Restaurant:     restaurant.Restaurant,
					DistanceMeters: restaurant.DistanceMeters,
				})
			}
		}
	}

	return results
}

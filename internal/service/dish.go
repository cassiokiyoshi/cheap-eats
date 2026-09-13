package service

import (
	"context"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
)

type DishService struct {
	dishRepository       repository.DishStore
	restaurantRepository repository.RestaurantStore
}

func NewDishService(
	dishRepository repository.DishStore,
	restaurantRepository repository.RestaurantStore,
) *DishService {
	return &DishService{
		dishRepository:       dishRepository,
		restaurantRepository: restaurantRepository,
	}
}

func (s *DishService) SearchNearby(
	ctx context.Context,
	latitude float64,
	longitude float64,
	radiusMeters float64,
	maxPrice int,
) ([]models.DishSearchResult, error) {
	restaurants, err := s.restaurantRepository.Nearby(
		ctx,
		latitude,
		longitude,
		radiusMeters,
	)
	if err != nil {
		return nil, err
	}

	dishes, err := s.dishRepository.ListByMaxPrice(
		ctx,
		maxPrice,
	)
	if err != nil {
		return nil, err
	}
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

	return results, nil
}

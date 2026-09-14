package service

import (
	"context"
	"fmt"
	"sort"

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
	sortBy string,
	limit int,
	offset int,
) ([]models.DishSearchResult, error) {
	if limit < 1 || limit > 100 || offset < 0 {
		return nil, fmt.Errorf("invalid pagination")
	}
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

	switch sortBy {
	case "price":
		sort.SliceStable(results, func(i, j int) bool {
			if results[i].Dish.Price == results[j].Dish.Price {
				return results[i].DistanceMeters <
					results[j].DistanceMeters
			}

			return results[i].Dish.Price <
				results[j].Dish.Price
		})

	default:
		sort.SliceStable(results, func(i, j int) bool {
			if results[i].DistanceMeters ==
				results[j].DistanceMeters {
				return results[i].Dish.Price <
					results[j].Dish.Price
			}

			return results[i].DistanceMeters <
				results[j].DistanceMeters
		})
	}

	if offset >= len(results) {
		return []models.DishSearchResult{}, nil
	}

	results = results[offset:]

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

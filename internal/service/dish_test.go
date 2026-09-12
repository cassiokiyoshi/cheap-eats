package service

import (
	"testing"

	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
)

func TestSearchNearbyFiltersByBudget(t *testing.T) {
	dishRepository := repository.NewDishRepository()
	restaurantRepository := repository.NewRestaurantRepository()

	service := NewDishService(
		dishRepository,
		restaurantRepository,
	)

	results := service.SearchNearby(
		35.6812,
		139.7671,
		2000,
		700,
	)

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	if results[0].Dish.Name != "Gyudon" {
		t.Errorf(
			"expected Gyudon, got %q",
			results[0].Dish.Name,
		)
	}

	if results[0].Dish.Price > 700 {
		t.Errorf(
			"expected price at most 700, got %d",
			results[0].Dish.Price,
		)
	}
}

func TestSearchNearbyFiltersByDistance(t *testing.T) {
	dishRepository := repository.NewDishRepository()
	restaurantRepository := repository.NewRestaurantRepository()

	service := NewDishService(
		dishRepository,
		restaurantRepository,
	)

	results := service.SearchNearby(
		35.6812,
		139.7671,
		500,
		1000,
	)

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	if results[0].Restaurant.ID != 1 {
		t.Errorf(
			"expected restaurant ID 1, got %d",
			results[0].Restaurant.ID,
		)
	}

	if results[0].Dish.Name != "Shoyu Ramen" {
		t.Errorf(
			"expected Shoyu Ramen, got %q",
			results[0].Dish.Name,
		)
	}
}

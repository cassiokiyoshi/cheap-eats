package repository

import (
	"context"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type DishStore interface {
	List() []models.Dish
	ListByMaxPrice(maxPrice int) []models.Dish
	FindByID(id int64) (models.Dish, bool)
	Create(dish models.Dish) models.Dish
}

type RestaurantStore interface {
	List(ctx context.Context) ([]models.Restaurant, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (models.Restaurant, bool, error)

	Nearby(
		ctx context.Context,
		latitude float64,
		longitude float64,
		radiusMeters float64,
	) ([]models.RestaurantSuggestion, error)
}

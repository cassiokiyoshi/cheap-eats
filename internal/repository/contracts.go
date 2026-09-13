package repository

import (
	"context"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type DishStore interface {
	List(ctx context.Context) ([]models.Dish, error)

	ListByMaxPrice(
		ctx context.Context,
		maxPrice int,
	) ([]models.Dish, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (models.Dish, bool, error)

	Create(
		ctx context.Context,
		dish models.Dish,
	) (models.Dish, error)
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

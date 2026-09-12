package repository

import "github.com/cassiokiyoshi/cheap-eats/internal/models"

type DishStore interface {
	List() []models.Dish
	ListByMaxPrice(maxPrice int) []models.Dish
	FindByID(id int64) (models.Dish, bool)
	Create(dish models.Dish) models.Dish
}

type RestaurantStore interface {
	List() []models.Restaurant
	FindByID(id int64) (models.Restaurant, bool)

	Nearby(
		latitude float64,
		longitude float64,
		radiusMeters float64,
	) []models.RestaurantSuggestion
}

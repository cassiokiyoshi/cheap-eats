package repository

import "github.com/cassiokiyoshi/cheap-eats/internal/models"

type DishRepository struct {
	dishes []models.Dish
}

func NewDishRepository() *DishRepository {
	return &DishRepository{
		dishes: []models.Dish{
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
		},
	}
}

func (r *DishRepository) List() []models.Dish {
	dishes := make([]models.Dish, len(r.dishes))
	copy(dishes, r.dishes)

	return dishes
}

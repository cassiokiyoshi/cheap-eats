package repository

import (
	"context"
	"sync"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
)

type DishRepository struct {
	mu     sync.RWMutex
	dishes []models.Dish
	nextID int64
}

func NewDishRepository() *DishRepository {
	return &DishRepository{
		dishes: []models.Dish{
			{
				ID:           1,
				RestaurantID: 1,
				Name:         "Shoyu Ramen",
				Price:        850,
				Currency:     "JPY",
			},
			{
				ID:           2,
				RestaurantID: 2,
				Name:         "Gyudon",
				Price:        650,
				Currency:     "JPY",
			},
		},
		nextID: 3,
	}
}

func (r *DishRepository) List(
	_ context.Context,
) ([]models.Dish, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dishes := make([]models.Dish, len(r.dishes))
	copy(dishes, r.dishes)

	return dishes, nil
}

func (r *DishRepository) ListByMaxPrice(
	_ context.Context,
	maxPrice int,
) ([]models.Dish, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dishes := make([]models.Dish, 0)

	for _, dish := range r.dishes {
		if dish.Price <= maxPrice {
			dishes = append(dishes, dish)
		}
	}

	return dishes, nil
}

func (r *DishRepository) FindByID(
	_ context.Context,
	id int64,
) (models.Dish, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dish := range r.dishes {
		if dish.ID == id {
			return dish, true, nil
		}
	}

	return models.Dish{}, false, nil
}

func (r *DishRepository) Create(
	_ context.Context,
	dish models.Dish,
) (models.Dish, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dish.ID = r.nextID
	r.nextID++

	r.dishes = append(r.dishes, dish)

	return dish, nil
}

func (r *DishRepository) ListByRestaurantID(
	_ context.Context,
	restaurantID int64,
) ([]models.Dish, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dishes := make([]models.Dish, 0)

	for _, dish := range r.dishes {
		if dish.RestaurantID == restaurantID {
			dishes = append(dishes, dish)
		}
	}

	return dishes, nil
}

package repository

import (
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

func (r *DishRepository) List() []models.Dish {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dishes := make([]models.Dish, len(r.dishes))
	copy(dishes, r.dishes)

	return dishes
}

func (r *DishRepository) ListByMaxPrice(
	maxPrice int,
) []models.Dish {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dishes := make([]models.Dish, 0)

	for _, dish := range r.dishes {
		if dish.Price <= maxPrice {
			dishes = append(dishes, dish)
		}
	}

	return dishes
}

func (r *DishRepository) FindByID(id int64) (models.Dish, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dish := range r.dishes {
		if dish.ID == id {
			return dish, true
		}
	}

	return models.Dish{}, false
}

func (r *DishRepository) Create(dish models.Dish) models.Dish {
	r.mu.Lock()
	defer r.mu.Unlock()

	dish.ID = r.nextID
	r.nextID++

	r.dishes = append(r.dishes, dish)

	return dish
}

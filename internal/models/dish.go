package models

type Dish struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Currency     string `json:"currency"`
	RestaurantID int64  `json:"restaurant_id"`
}

package models

type Dish struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	NameJA       *string `json:"name_ja"`
	NameEN       *string `json:"name_en"`
	Price        int     `json:"price"`
	Currency     string  `json:"currency"`
	RestaurantID int64   `json:"restaurant_id"`
}

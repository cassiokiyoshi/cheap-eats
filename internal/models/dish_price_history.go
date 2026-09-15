package models

import "time"

type DishPriceHistory struct {
	ID        int64     `json:"id"`
	DishID    int64     `json:"dish_id"`
	OldPrice  int       `json:"old_price"`
	NewPrice  int       `json:"new_price"`
	Currency  string    `json:"currency"`
	ChangedAt time.Time `json:"changed_at"`
}

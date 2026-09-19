package models

type Restaurant struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	NameJA    *string `json:"name_ja"`
	NameEN    *string `json:"name_en"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RestaurantSuggestion struct {
	Restaurant
	DistanceMeters int `json:"distance_meters"`
}

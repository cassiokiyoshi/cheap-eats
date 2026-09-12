package models

type DishSearchResult struct {
	Dish           Dish       `json:"dish"`
	Restaurant     Restaurant `json:"restaurant"`
	DistanceMeters int        `json:"distance_meters"`
}

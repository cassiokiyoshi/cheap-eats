package repository

import "testing"

func TestRestaurantRepositoryNearby(t *testing.T) {
	repository := NewRestaurantRepository()

	restaurants := repository.Nearby(
		35.6812,
		139.7671,
		500,
	)

	if len(restaurants) != 1 {
		t.Fatalf(
			"expected 1 nearby restaurant, got %d",
			len(restaurants),
		)
	}

	if restaurants[0].ID != 1 {
		t.Errorf(
			"expected restaurant ID 1, got %d",
			restaurants[0].ID,
		)
	}

	if restaurants[0].DistanceMeters != 0 {
		t.Errorf(
			"expected distance 0, got %d",
			restaurants[0].DistanceMeters,
		)
	}
}

func TestRestaurantRepositoryNearbyOrdersByDistance(t *testing.T) {
	repository := NewRestaurantRepository()

	restaurants := repository.Nearby(
		35.6812,
		139.7671,
		2000,
	)

	if len(restaurants) != 2 {
		t.Fatalf(
			"expected 2 nearby restaurants, got %d",
			len(restaurants),
		)
	}

	if restaurants[0].ID != 1 {
		t.Errorf(
			"expected closest restaurant ID 1, got %d",
			restaurants[0].ID,
		)
	}

	if restaurants[1].ID != 2 {
		t.Errorf(
			"expected second restaurant ID 2, got %d",
			restaurants[1].ID,
		)
	}

	if restaurants[0].DistanceMeters > restaurants[1].DistanceMeters {
		t.Error("expected restaurants to be ordered by distance")
	}
}

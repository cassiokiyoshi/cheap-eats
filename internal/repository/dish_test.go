package repository

import "testing"

func TestDishRepositoryFindByID(t *testing.T) {
	repository := NewDishRepository()

	dish, found := repository.FindByID(1)

	if !found {
		t.Fatal("expected dish to be found")
	}

	if dish.ID != 1 {
		t.Errorf("expected dish ID 1, got %d", dish.ID)
	}

	if dish.Name != "Shoyu Ramen" {
		t.Errorf("expected Shoyu Ramen, got %s", dish.Name)
	}
}

func TestDishRepositoryFindByIDNotFound(t *testing.T) {
	repository := NewDishRepository()

	_, found := repository.FindByID(999)

	if found {
		t.Error("expected dish not to be found")
	}
}

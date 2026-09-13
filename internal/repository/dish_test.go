package repository

import (
	"context"
	"testing"
)

func TestDishRepositoryFindByID(t *testing.T) {
	repository := NewDishRepository()

	dish, found, err := repository.FindByID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Fatalf("find dish: %v", err)
	}

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

	_, found, err := repository.FindByID(
		context.Background(),
		999,
	)

	if err != nil {
		t.Fatalf("find dish: %v", err)
	}

	if found {
		t.Error("expected dish not to be found")
	}
}

func TestDishRepositoryListByMaxPrice(t *testing.T) {
	repository := NewDishRepository()

	dishes, err := repository.ListByMaxPrice(
		context.Background(),
		700,
	)
	if err != nil {
		t.Fatalf("list dishes by maximum price: %v", err)
	}

	if len(dishes) != 1 {
		t.Fatalf(
			"expected 1 dish, got %d",
			len(dishes),
		)
	}

	if dishes[0].Name != "Gyudon" {
		t.Errorf(
			"expected Gyudon, got %q",
			dishes[0].Name,
		)
	}

	if dishes[0].Price > 700 {
		t.Errorf(
			"expected price at most 700, got %d",
			dishes[0].Price,
		)
	}
}

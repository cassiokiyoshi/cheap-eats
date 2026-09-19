package repository

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresBilingualNames(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").
		Scan(&databaseName); err != nil {
		t.Fatal(err)
	}
	if databaseName != "cheap_eats_test" {
		t.Fatalf("expected cheap_eats_test, got %q", databaseName)
	}

	restaurants := NewPostgresRestaurantRepository(pool)
	dishes := NewPostgresDishRepository(pool)

	restaurantJA := "さくら食堂"
	restaurantEN := "Sakura Shokudo"
	dishJA := "醤油ラーメン"
	dishEN := "Shoyu Ramen"

	for _, bilingual := range []bool{false, true} {
		name := "null names"
		if bilingual {
			name = "bilingual names"
		}

		t.Run(name, func(t *testing.T) {
			input := models.Restaurant{
				Name:      "Bilingual Test Restaurant",
				Address:   "Test address",
				Latitude:  0,
				Longitude: 0,
			}
			if bilingual {
				input.NameJA = &restaurantJA
				input.NameEN = &restaurantEN
			}

			restaurant, err := restaurants.Create(ctx, input)
			if err != nil {
				t.Fatalf("create restaurant: %v", err)
			}

			// Delete only this test's restaurant and its cascading dishes.
			defer func() {
				cleanupCtx, cleanupCancel := context.WithTimeout(
					context.Background(), 5*time.Second,
				)
				defer cleanupCancel()

				_, err := pool.Exec(
					cleanupCtx,
					"DELETE FROM restaurants WHERE id = $1",
					restaurant.ID,
				)
				if err != nil {
					t.Errorf("cleanup: %v", err)
				}
			}()

			savedRestaurant, found, err := restaurants.FindByID(
				ctx, restaurant.ID,
			)
			if err != nil || !found {
				t.Fatalf("find restaurant: found=%v, err=%v", found, err)
			}
			input.ID = restaurant.ID
			if !reflect.DeepEqual(savedRestaurant, input) {
				t.Fatal("restaurant did not round-trip correctly")
			}

			dishInput := models.Dish{
				RestaurantID: restaurant.ID,
				Name:         "Bilingual Test Ramen",
				Price:        850,
				Currency:     "JPY",
			}
			if bilingual {
				dishInput.NameJA = &dishJA
				dishInput.NameEN = &dishEN
			}

			dish, err := dishes.Create(ctx, dishInput)
			if err != nil {
				t.Fatalf("create dish: %v", err)
			}

			savedDish, found, err := dishes.FindByID(ctx, dish.ID)
			if err != nil || !found {
				t.Fatalf("find dish: found=%v, err=%v", found, err)
			}
			dishInput.ID = dish.ID
			if !reflect.DeepEqual(savedDish, dishInput) {
				t.Fatal("dish did not round-trip correctly")
			}

			listed, err := dishes.ListByRestaurantID(ctx, restaurant.ID)
			if err != nil {
				t.Fatal(err)
			}
			if len(listed) != 1 || !reflect.DeepEqual(listed[0], dishInput) {
				t.Fatal("restaurant dish list did not preserve names")
			}

			results, err := dishes.SearchNearby(
				ctx, 0, 0, 1, 1000, "distance", 100, 0,
			)
			if err != nil {
				t.Fatalf("nearby search: %v", err)
			}

			matched := false
			for _, result := range results {
				if result.Dish.ID != dish.ID {
					continue
				}
				matched = true

				if !reflect.DeepEqual(result.Dish, dishInput) {
					t.Error("nearby search did not preserve dish names")
				}
				if !reflect.DeepEqual(result.Restaurant, input) {
					t.Error("nearby search did not preserve restaurant names")
				}
			}
			if !matched {
				t.Fatal("nearby search did not return the test dish")
			}

			updated, found, err := dishes.UpdatePrice(ctx, dish.ID, 900)
			if err != nil || !found {
				t.Fatalf("update price: found=%v, err=%v", found, err)
			}
			dishInput.Price = 900
			if !reflect.DeepEqual(updated, dishInput) {
				t.Fatal("price update did not preserve names")
			}
		})
	}
}

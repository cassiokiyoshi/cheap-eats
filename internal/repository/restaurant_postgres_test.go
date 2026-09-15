package repository

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRestaurantCreate(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create test database pool: %v", err)
	}
	defer pool.Close()

	var databaseName string
	err = pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}

	if databaseName != "cheap_eats_test" {
		t.Fatalf("expected cheap_eats_test database, got %q", databaseName)
	}

	store := NewPostgresRestaurantRepository(pool)

	input := models.Restaurant{
		Name:      "Integration Test Curry",
		Address:   "Test address, Tokyo",
		Latitude:  35.6830,
		Longitude: 139.7690,
	}

	created, err := store.Create(ctx, input)
	if err != nil {
		t.Fatalf("create restaurant: %v", err)
	}

	// Registered after pool.Close, so this runs before the pool closes.
	defer func() {
		cleanupContext, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := pool.Exec(
			cleanupContext,
			"DELETE FROM restaurants WHERE id = $1",
			created.ID,
		)
		if err != nil {
			t.Errorf("clean up test restaurant: %v", err)
		}
	}()

	if created.ID <= 0 {
		t.Fatalf("expected a generated ID, got %d", created.ID)
	}

	saved, found, err := store.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("read restaurant: %v", err)
	}
	if !found {
		t.Fatal("created restaurant was not found")
	}

	if saved.Name != input.Name || saved.Address != input.Address {
		t.Errorf("unexpected saved restaurant: %+v", saved)
	}

	const tolerance = 0.000001

	if math.Abs(saved.Latitude-input.Latitude) > tolerance {
		t.Errorf(
			"expected latitude %f, got %f",
			input.Latitude,
			saved.Latitude,
		)
	}

	if math.Abs(saved.Longitude-input.Longitude) > tolerance {
		t.Errorf(
			"expected longitude %f, got %f",
			input.Longitude,
			saved.Longitude,
		)
	}

	t.Run("list restaurant dishes", func(t *testing.T) {
		dishStore := NewPostgresDishRepository(pool)

		// A newly created restaurant should have no dishes.
		empty, err := dishStore.ListByRestaurantID(ctx, created.ID)
		if err != nil {
			t.Fatalf("list empty restaurant dishes: %v", err)
		}
		if empty == nil || len(empty) != 0 {
			t.Fatalf("expected a non-nil empty slice, got %#v", empty)
		}

		first, err := dishStore.Create(ctx, models.Dish{
			RestaurantID: created.ID,
			Name:         "Integration Test Curry",
			Price:        900,
			Currency:     "JPY",
		})
		if err != nil {
			t.Fatalf("create first dish: %v", err)
		}

		second, err := dishStore.Create(ctx, models.Dish{
			RestaurantID: created.ID,
			Name:         "Integration Test Soup",
			Price:        500,
			Currency:     "JPY",
		})
		if err != nil {
			t.Fatalf("create second dish: %v", err)
		}

		dishes, err := dishStore.ListByRestaurantID(ctx, created.ID)
		if err != nil {
			t.Fatalf("list restaurant dishes: %v", err)
		}

		if len(dishes) != 2 {
			t.Fatalf("expected 2 dishes, got %d", len(dishes))
		}

		if dishes[0] != first || dishes[1] != second {
			t.Errorf(
				"expected dishes in ID order: %+v, %+v; got %+v",
				first,
				second,
				dishes,
			)
		}
	})

	t.Run("update dish price", func(t *testing.T) {
		dishStore := NewPostgresDishRepository(pool)

		original, err := dishStore.Create(ctx, models.Dish{
			RestaurantID: created.ID,
			Name:         "Price Update Test Curry",
			Price:        900,
			Currency:     "JPY",
		})
		if err != nil {
			t.Fatalf("create test dish: %v", err)
		}

		// Set a known old timestamp so the test needs no sleep.
		oldTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

		_, err = pool.Exec(
			ctx,
			"UPDATE dishes SET updated_at = $2 WHERE id = $1",
			original.ID,
			oldTime,
		)
		if err != nil {
			t.Fatalf("set initial timestamp: %v", err)
		}

		updated, found, err := dishStore.UpdatePrice(
			ctx,
			original.ID,
			950,
		)
		if err != nil {
			t.Fatalf("update price: %v", err)
		}
		if !found {
			t.Fatal("test dish was not found")
		}

		want := original
		want.Price = 950

		if updated != want {
			t.Errorf("expected updated dish %+v, got %+v", want, updated)
		}

		saved, found, err := dishStore.FindByID(ctx, original.ID)
		if err != nil {
			t.Fatalf("read updated dish: %v", err)
		}
		if !found || saved != want {
			t.Errorf(
				"price was not saved correctly: found=%v, got=%+v",
				found,
				saved,
			)
		}

		var updatedAt time.Time
		err = pool.QueryRow(
			ctx,
			"SELECT updated_at FROM dishes WHERE id = $1",
			original.ID,
		).Scan(&updatedAt)
		if err != nil {
			t.Fatalf("read timestamp: %v", err)
		}
		if !updatedAt.After(oldTime) {
			t.Errorf("updated_at did not advance: got %v", updatedAt)
		}

		// Verify the recorded price change.
		var oldPrice, newPrice int
		var currency string
		var changedAt time.Time

		err = pool.QueryRow(ctx, `
			SELECT old_price, new_price, currency, changed_at
			FROM dish_price_history
			WHERE dish_id = $1
			ORDER BY id DESC
			LIMIT 1
		`, original.ID).Scan(
			&oldPrice,
			&newPrice,
			&currency,
			&changedAt,
		)
		if err != nil {
			t.Fatalf("read price history: %v", err)
		}

		if oldPrice != 900 || newPrice != 950 || currency != "JPY" {
			t.Errorf(
				"unexpected history: %d -> %d %s",
				oldPrice,
				newPrice,
				currency,
			)
		}

		if !changedAt.Equal(updatedAt) {
			t.Errorf(
				"history timestamp %v differs from dish timestamp %v",
				changedAt,
				updatedAt,
			)
		}

		// Repeating the same price should not create another change.
		repeated, found, err := dishStore.UpdatePrice(
			ctx,
			original.ID,
			950,
		)
		if err != nil {
			t.Fatalf("repeat price update: %v", err)
		}
		if !found || repeated != want {
			t.Errorf(
				"unexpected repeated update: found=%v, dish=%+v",
				found,
				repeated,
			)
		}

		var historyCount int
		err = pool.QueryRow(
			ctx,
			"SELECT COUNT(*) FROM dish_price_history WHERE dish_id = $1",
			original.ID,
		).Scan(&historyCount)
		if err != nil {
			t.Fatalf("count price history: %v", err)
		}
		if historyCount != 1 {
			t.Errorf("expected 1 history entry, got %d", historyCount)
		}

		var repeatedUpdatedAt time.Time
		err = pool.QueryRow(
			ctx,
			"SELECT updated_at FROM dishes WHERE id = $1",
			original.ID,
		).Scan(&repeatedUpdatedAt)
		if err != nil {
			t.Fatalf("read timestamp after repeated update: %v", err)
		}
		if !repeatedUpdatedAt.Equal(updatedAt) {
			t.Error("same-price request changed updated_at")
		}

		// Delete only our test dish to obtain a known missing ID.
		_, err = pool.Exec(
			ctx,
			"DELETE FROM dishes WHERE id = $1",
			original.ID,
		)
		if err != nil {
			t.Fatalf("delete test dish: %v", err)
		}

		_, found, err = dishStore.UpdatePrice(ctx, original.ID, 1000)
		if err != nil {
			t.Fatalf("update missing dish: %v", err)
		}
		if found {
			t.Error("expected found=false for a deleted dish")
		}
	})

	t.Run("list price history", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		dishStore := NewPostgresDishRepository(pool)

		first, err := dishStore.Create(ctx, models.Dish{
			RestaurantID: created.ID,
			Name:         "History Test Curry",
			Price:        800,
			Currency:     "JPY",
		})
		if err != nil {
			t.Fatalf("create first dish: %v", err)
		}

		second, err := dishStore.Create(ctx, models.Dish{
			RestaurantID: created.ID,
			Name:         "History Test Soup",
			Price:        500,
			Currency:     "JPY",
		})
		if err != nil {
			t.Fatalf("create second dish: %v", err)
		}

		empty, err := dishStore.ListPriceHistory(ctx, first.ID)
		if err != nil {
			t.Fatalf("read empty history: %v", err)
		}
		if empty == nil || len(empty) != 0 {
			t.Fatalf("expected non-nil empty history, got %#v", empty)
		}

		updates := []struct {
			id    int64
			price int
		}{
			{first.ID, 850},
			{second.ID, 550},
			{first.ID, 900},
			{first.ID, 900},
		}

		for _, update := range updates {
			_, found, err := dishStore.UpdatePrice(
				ctx,
				update.id,
				update.price,
			)
			if err != nil || !found {
				t.Fatalf("prepare history: found=%v, err=%v", found, err)
			}
		}

		history, err := dishStore.ListPriceHistory(ctx, first.ID)
		if err != nil {
			t.Fatalf("list price history: %v", err)
		}
		if len(history) != 2 {
			t.Fatalf("expected 2 history entries, got %d", len(history))
		}

		wantChanges := [][2]int{
			{850, 900},
			{800, 850},
		}

		for i, entry := range history {
			if entry.DishID != first.ID || entry.Currency != "JPY" {
				t.Errorf("unexpected history entry: %+v", entry)
			}
			if entry.OldPrice != wantChanges[i][0] ||
				entry.NewPrice != wantChanges[i][1] {
				t.Errorf("unexpected change at index %d: %+v", i, entry)
			}
			if entry.ID <= 0 || entry.ChangedAt.IsZero() {
				t.Errorf("missing ID or timestamp: %+v", entry)
			}
		}

		if history[0].ID <= history[1].ID {
			t.Error("expected history in descending ID order")
		}

		otherHistory, err := dishStore.ListPriceHistory(ctx, second.ID)
		if err != nil {
			t.Fatalf("list second dish history: %v", err)
		}
		if len(otherHistory) != 1 {
			t.Fatalf("expected 1 entry for second dish, got %d", len(otherHistory))
		}
		if otherHistory[0].DishID != second.ID ||
			otherHistory[0].OldPrice != 500 ||
			otherHistory[0].NewPrice != 550 {
			t.Errorf("unexpected second dish history: %+v", otherHistory[0])
		}
	})
}

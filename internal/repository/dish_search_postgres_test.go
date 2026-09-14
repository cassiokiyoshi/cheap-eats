package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresDishSearchNearby(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	repository := NewPostgresDishRepository(pool)

	tests := []struct {
		name     string
		radius   float64
		maxPrice int
		sortBy   string
		limit    int
		offset   int
		want     []string
	}{
		{
			name:     "nearest first",
			radius:   2000,
			maxPrice: 1000,
			sortBy:   "distance",
			want:     []string{"Shoyu Ramen", "Gyudon"},
		},
		{
			name:     "cheapest first",
			radius:   2000,
			maxPrice: 1000,
			sortBy:   "price",
			want:     []string{"Gyudon", "Shoyu Ramen"},
		},
		{
			name:     "budget includes exact price",
			radius:   2000,
			maxPrice: 650,
			sortBy:   "price",
			want:     []string{"Gyudon"},
		},
		{
			name:     "radius excludes distant restaurant",
			radius:   500,
			maxPrice: 1000,
			sortBy:   "distance",
			want:     []string{"Shoyu Ramen"},
		},
		{
			name:     "both filters must match",
			radius:   500,
			maxPrice: 700,
			sortBy:   "distance",
			want:     []string{},
		},
		{
			name:     "first page by price",
			radius:   2000,
			maxPrice: 1000,
			sortBy:   "price",
			limit:    1,
			offset:   0,
			want:     []string{"Gyudon"},
		},
		{
			name:     "second page by price",
			radius:   2000,
			maxPrice: 1000,
			sortBy:   "price",
			limit:    1,
			offset:   1,
			want:     []string{"Shoyu Ramen"},
		},
		{
			name:     "page beyond results",
			radius:   2000,
			maxPrice: 1000,
			sortBy:   "price",
			limit:    1,
			offset:   2,
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.limit
			if limit == 0 {
				limit = 20
			}
			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			results, err := repository.SearchNearby(
				ctx,
				35.6812,
				139.7671,
				tt.radius,
				tt.maxPrice,
				tt.sortBy,
				limit,
				tt.offset,
			)
			if err != nil {
				t.Fatalf("search nearby dishes: %v", err)
			}

			if results == nil {
				t.Fatal("expected a non-nil result slice, even when empty")
			}

			if len(results) != len(tt.want) {
				t.Fatalf(
					"expected %d results, got %d",
					len(tt.want),
					len(results),
				)
			}

			for i, result := range results {
				if result.Dish.Name != tt.want[i] {
					t.Errorf(
						"result %d: expected %q, got %q",
						i,
						tt.want[i],
						result.Dish.Name,
					)
				}

				if result.Dish.RestaurantID != result.Restaurant.ID {
					t.Errorf("result %d has a mismatched restaurant", i)
				}
			}
		})
	}
}

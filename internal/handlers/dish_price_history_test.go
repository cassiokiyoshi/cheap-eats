package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestDishPriceHistory(t *testing.T) {
	dishes := repository.NewDishRepository()
	handler := NewDishHandler(
		dishes,
		repository.NewRestaurantRepository(),
	)

	updates := []struct {
		dishID int64
		price  int
	}{
		{1, 900},
		{2, 700}, // Another dish: must not appear in dish 1's history.
		{1, 950},
		{1, 950}, // Same price: must not add another entry.
	}

	for _, update := range updates {
		_, found, err := dishes.UpdatePrice(
			context.Background(),
			update.dishID,
			update.price,
		)
		if err != nil || !found {
			t.Fatalf("prepare history: found=%v, err=%v", found, err)
		}
	}

	router := chi.NewRouter()
	router.Get("/api/dishes/{id}/price-history", handler.PriceHistory)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/dishes/1/price-history",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected application/json, got %q", got)
	}

	var history []models.DishPriceHistory
	if err := json.Unmarshal(response.Body.Bytes(), &history); err != nil {
		t.Fatalf("decode history: %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(history))
	}

	wantChanges := [][2]int{
		{900, 950},
		{850, 900},
	}

	for i, entry := range history {
		if entry.DishID != 1 || entry.Currency != "JPY" {
			t.Errorf("unexpected entry: %+v", entry)
		}
		if entry.OldPrice != wantChanges[i][0] ||
			entry.NewPrice != wantChanges[i][1] {
			t.Errorf("unexpected price change at index %d: %+v", i, entry)
		}
		if entry.ID <= 0 || entry.ChangedAt.IsZero() {
			t.Errorf("missing ID or timestamp: %+v", entry)
		}
	}

	if history[0].ID <= history[1].ID {
		t.Error("expected history in descending ID order")
	}
}

func TestDishPriceHistoryEmptyAndInvalidRequests(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantBody   string
	}{
		{"no history", "1", 200, "[]"},
		{"missing dish", "999", 404, "dish not found"},
		{"invalid ID", "abc", 400, "invalid dish ID"},
		{"zero ID", "0", 400, "invalid dish ID"},
		{"negative ID", "-1", 400, "invalid dish ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewDishHandler(
				repository.NewDishRepository(),
				repository.NewRestaurantRepository(),
			)

			router := chi.NewRouter()
			router.Get(
				"/api/dishes/{id}/price-history",
				handler.PriceHistory,
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/dishes/"+tt.id+"/price-history",
				nil,
			)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					response.Code,
				)
			}

			if got := strings.TrimSpace(response.Body.String()); got != tt.wantBody {
				t.Errorf("expected body %q, got %q", tt.wantBody, got)
			}
		})
	}
}

func TestDishPriceHistoryPagination(t *testing.T) {
	dishes := repository.NewDishRepository()

	// Create two changes for dish 1, with another dish's change between them.
	updates := []struct {
		id    int64
		price int
	}{
		{1, 900},
		{2, 700},
		{1, 950},
	}

	for _, update := range updates {
		_, found, err := dishes.UpdatePrice(
			context.Background(),
			update.id,
			update.price,
		)
		if err != nil || !found {
			t.Fatalf("prepare history: found=%v, err=%v", found, err)
		}
	}

	handler := NewDishHandler(
		dishes,
		repository.NewRestaurantRepository(),
	)

	router := chi.NewRouter()
	router.Get("/api/dishes/{id}/price-history", handler.PriceHistory)

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantPrices []int
	}{
		{"defaults", "", 200, []int{950, 900}},
		{"first page", "?limit=1&offset=0", 200, []int{950}},
		{"second page", "?limit=1&offset=1", 200, []int{900}},
		{"past last page", "?limit=1&offset=2", 200, []int{}},
		{"maximum limit", "?limit=100", 200, []int{950, 900}},
		{"offset only", "?offset=1", 200, []int{900}},
		{"zero limit", "?limit=0", 400, nil},
		{"negative limit", "?limit=-1", 400, nil},
		{"limit too large", "?limit=101", 400, nil},
		{"empty limit", "?limit=", 400, nil},
		{"invalid limit", "?limit=abc", 400, nil},
		{"decimal limit", "?limit=1.5", 400, nil},
		{"negative offset", "?offset=-1", 400, nil},
		{"empty offset", "?offset=", 400, nil},
		{"invalid offset", "?offset=abc", 400, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodGet,
				"/api/dishes/1/price-history"+tt.query,
				nil,
			)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d: %s",
					tt.wantStatus,
					response.Code,
					response.Body.String(),
				)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var history []models.DishPriceHistory
			if err := json.Unmarshal(response.Body.Bytes(), &history); err != nil {
				t.Fatalf("decode history: %v", err)
			}

			if history == nil {
				t.Fatal("expected a JSON array, got null")
			}
			if len(history) != len(tt.wantPrices) {
				t.Fatalf(
					"expected %d entries, got %d",
					len(tt.wantPrices),
					len(history),
				)
			}

			for i, entry := range history {
				if entry.DishID != 1 || entry.NewPrice != tt.wantPrices[i] {
					t.Errorf("unexpected entry at index %d: %+v", i, entry)
				}
			}
		})
	}
}

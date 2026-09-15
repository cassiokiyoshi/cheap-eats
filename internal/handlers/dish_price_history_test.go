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

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestListDishesByRestaurant(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantDishID int64
		wantError  string
	}{
		{"first restaurant", "1", 200, 1, ""},
		{"second restaurant", "2", 200, 2, ""},
		{"missing restaurant", "999", 404, 0, "restaurant not found"},
		{"invalid ID", "abc", 400, 0, "invalid restaurant ID"},
		{"zero ID", "0", 400, 0, "invalid restaurant ID"},
		{"negative ID", "-1", 400, 0, "invalid restaurant ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dishes := repository.NewDishRepository()
			restaurants := repository.NewRestaurantRepository()
			handler := NewDishHandler(dishes, restaurants)

			router := chi.NewRouter()
			router.Get(
				"/api/restaurants/{id}/dishes",
				handler.ListByRestaurant,
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/restaurants/"+tt.id+"/dishes",
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

			if tt.wantError != "" {
				if got := strings.TrimSpace(response.Body.String()); got != tt.wantError {
					t.Errorf("expected error %q, got %q", tt.wantError, got)
				}
				return
			}

			var result []models.Dish
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if len(result) != 1 {
				t.Fatalf("expected one dish, got %d", len(result))
			}

			if result[0].ID != tt.wantDishID {
				t.Errorf("expected dish ID %d, got %d", tt.wantDishID, result[0].ID)
			}

			if got := strconv.FormatInt(result[0].RestaurantID, 10); got != tt.id {
				t.Errorf("dish belongs to restaurant %s, expected %s", got, tt.id)
			}
		})
	}
}

func TestListDishesByRestaurantReturnsEmptyArray(t *testing.T) {
	dishes := repository.NewDishRepository()
	restaurants := repository.NewRestaurantRepository()

	created, err := restaurants.Create(
		context.Background(),
		models.Restaurant{
			Name:      "Restaurant Without Dishes",
			Address:   "Tokyo",
			Latitude:  35.6830,
			Longitude: 139.7690,
		},
	)
	if err != nil {
		t.Fatalf("create restaurant: %v", err)
	}

	handler := NewDishHandler(dishes, restaurants)

	router := chi.NewRouter()
	router.Get(
		"/api/restaurants/{id}/dishes",
		handler.ListByRestaurant,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/restaurants/"+strconv.FormatInt(created.ID, 10)+"/dishes",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	if got := strings.TrimSpace(response.Body.String()); got != "[]" {
		t.Errorf("expected [], got %q", got)
	}
}

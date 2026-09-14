package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cassiokiyoshi/cheap-eats/internal/models"
	"github.com/cassiokiyoshi/cheap-eats/internal/repository"
)

func TestCreateRestaurant(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{"Tokyo coordinates", 35.6830, 139.7690},
		{"zero coordinates", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := repository.NewRestaurantRepository()
			handler := NewRestaurantHandler(store)

			body := fmt.Sprintf(
				`{
					"name": "  Curry Kitchen  ",
					"address": "  Tokyo  ",
					"latitude": %f,
					"longitude": %f
				}`,
				tt.latitude,
				tt.longitude,
			)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/restaurants",
				strings.NewReader(body),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler.Create(response, request)

			if response.Code != http.StatusCreated {
				t.Fatalf(
					"expected 201, got %d: %s",
					response.Code,
					response.Body.String(),
				)
			}

			var created models.Restaurant
			if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			want := models.Restaurant{
				ID:        3,
				Name:      "Curry Kitchen",
				Address:   "Tokyo",
				Latitude:  tt.latitude,
				Longitude: tt.longitude,
			}
			if created != want {
				t.Errorf("expected %+v, got %+v", want, created)
			}

			if got := response.Header().Get("Location"); got != "/api/restaurants/3" {
				t.Errorf("unexpected Location header: %q", got)
			}

			saved, found, err := store.FindByID(context.Background(), created.ID)
			if err != nil {
				t.Fatalf("find saved restaurant: %v", err)
			}
			if !found || saved != want {
				t.Errorf("restaurant not saved correctly: found=%v, got=%+v", found, saved)
			}
		})
	}
}

func TestCreateRestaurantRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			"blank name",
			`{"name":" ","address":"Tokyo","latitude":35,"longitude":139}`,
		},
		{
			"missing address",
			`{"name":"Curry","latitude":35,"longitude":139}`,
		},
		{
			"missing latitude",
			`{"name":"Curry","address":"Tokyo","longitude":139}`,
		},
		{
			"null longitude",
			`{"name":"Curry","address":"Tokyo","latitude":35,"longitude":null}`,
		},
		{
			"latitude too high",
			`{"name":"Curry","address":"Tokyo","latitude":91,"longitude":139}`,
		},
		{
			"latitude too low",
			`{"name":"Curry","address":"Tokyo","latitude":-91,"longitude":139}`,
		},
		{
			"longitude too high",
			`{"name":"Curry","address":"Tokyo","latitude":35,"longitude":181}`,
		},
		{
			"longitude too low",
			`{"name":"Curry","address":"Tokyo","latitude":35,"longitude":-181}`,
		},
		{
			"unknown field",
			`{"name":"Curry","address":"Tokyo","latitude":35,"longitude":139,"extra":true}`,
		},
		{
			"multiple JSON values",
			`{"name":"Curry","address":"Tokyo","latitude":35,"longitude":139} {}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := repository.NewRestaurantRepository()
			handler := NewRestaurantHandler(store)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/restaurants",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler.Create(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf(
					"expected 400, got %d: %s",
					response.Code,
					response.Body.String(),
				)
			}

			restaurants, err := store.List(context.Background())
			if err != nil {
				t.Fatalf("list restaurants: %v", err)
			}
			if len(restaurants) != 2 {
				t.Errorf(
					"invalid input changed the repository: got %d restaurants",
					len(restaurants),
				)
			}
		})
	}
}

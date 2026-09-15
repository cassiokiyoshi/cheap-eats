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
)

func TestCreateDish(t *testing.T) {
	dishRepository := repository.NewDishRepository()
	restaurantRepository := repository.NewRestaurantRepository()
	dishHandler := NewDishHandler(dishRepository, restaurantRepository)

	body := `{
		"restaurant_id": 1,
		"name": "Chicken Curry",
		"price": 900,
		"currency": "jpy"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/dishes",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	dishHandler.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.Code,
		)
	}

	var dish models.Dish

	if err := json.NewDecoder(response.Body).Decode(&dish); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if dish.ID != 3 {
		t.Errorf("expected ID 3, got %d", dish.ID)
	}

	if dish.Name != "Chicken Curry" {
		t.Errorf("expected Chicken Curry, got %q", dish.Name)
	}

	if dish.Currency != "JPY" {
		t.Errorf("expected JPY, got %q", dish.Currency)
	}

	if dish.RestaurantID != 1 {
		t.Errorf("expected restaurant ID 1, got %d", dish.RestaurantID)
	}
}

func TestCreateDishRejectsInvalidBody(t *testing.T) {
	valid := `{"restaurant_id":1,"name":"Curry","price":900,"currency":"JPY"}`

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "empty body",
			body:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "two JSON values",
			body:   valid + ` {}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "trailing garbage",
			body:   valid + ` invalid`,
			status: http.StatusBadRequest,
		},
		{
			name:   "unknown field",
			body:   `{"restaurant_id":1,"name":"Curry","price":900,"extra":true}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "oversized body",
			body:   valid + strings.Repeat(" ", 64*1024),
			status: http.StatusRequestEntityTooLarge,
		},
		{
			name:   "price exceeds database limit",
			body:   `{"restaurant_id":1,"name":"Curry","price":2147483648,"currency":"JPY"}`,
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dishes := repository.NewDishRepository()
			restaurants := repository.NewRestaurantRepository()
			handler := NewDishHandler(dishes, restaurants)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/dishes",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler.Create(response, request)

			if response.Code != tt.status {
				t.Fatalf(
					"expected status %d, got %d: %s",
					tt.status,
					response.Code,
					response.Body.String(),
				)
			}

			saved, err := dishes.List(context.Background())
			if err != nil {
				t.Fatalf("list dishes: %v", err)
			}
			if len(saved) != 2 {
				t.Errorf("rejected request changed dish count to %d", len(saved))
			}
		})
	}
}

func TestCreateDishCurrency(t *testing.T) {
	tests := []struct {
		name       string
		currency   string
		omit       bool
		wantStatus int
	}{
		{"JPY", "JPY", false, 201},
		{"lowercase", "jpy", false, 201},
		{"surrounding spaces", " jpy ", false, 201},
		{"blank defaults to JPY", "", false, 201},
		{"omitted defaults to JPY", "", true, 201},
		{"unsupported currency", "USD", false, 400},
		{"misspelled currency", "JYP", false, 400},
		{"too short", "JP", false, 400},
		{"too long", "JPYY", false, 400},
		{"numbers", "123", false, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dishes := repository.NewDishRepository()
			handler := NewDishHandler(
				dishes,
				repository.NewRestaurantRepository(),
			)

			input := map[string]any{
				"restaurant_id": 1,
				"name":          "Curry",
				"price":         900,
			}
			if !tt.omit {
				input["currency"] = tt.currency
			}

			body, err := json.Marshal(input)
			if err != nil {
				t.Fatalf("encode request: %v", err)
			}

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/dishes",
				strings.NewReader(string(body)),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler.Create(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d: %s",
					tt.wantStatus,
					response.Code,
					response.Body.String(),
				)
			}

			saved, err := dishes.List(context.Background())
			if err != nil {
				t.Fatalf("list dishes: %v", err)
			}

			if tt.wantStatus == http.StatusBadRequest {
				if strings.TrimSpace(response.Body.String()) != "currency must be JPY" {
					t.Errorf("unexpected error: %s", response.Body.String())
				}
				if len(saved) != 2 {
					t.Error("invalid currency caused a dish to be saved")
				}
				return
			}

			var created models.Dish
			if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if created.Currency != "JPY" {
				t.Errorf("expected JPY, got %q", created.Currency)
			}

			stored, found, err := dishes.FindByID(context.Background(), created.ID)
			if err != nil || !found || stored.Currency != "JPY" {
				t.Fatalf(
					"JPY was not saved correctly: found=%v, dish=%+v, err=%v",
					found,
					stored,
					err,
				)
			}
		})
	}
}

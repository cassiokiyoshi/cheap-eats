package handlers

import (
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
	dishHandler := NewDishHandler(dishRepository)

	body := `{
		"name": "Chicken Curry",
		"price": 900,
		"currency": "jpy",
		"restaurant_name": "Curry House"
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
}

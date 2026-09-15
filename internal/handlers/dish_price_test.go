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

func TestUpdateDishPrice(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		wantStatus int
		wantPrice  int
	}{
		{"valid price", "1", `{"price":950}`, 200, 950},
		{"minimum price", "1", `{"price":1}`, 200, 1},
		{"maximum price", "1", `{"price":2147483647}`, 200, 2147483647},
		{"zero price", "1", `{"price":0}`, 400, 850},
		{"negative price", "1", `{"price":-10}`, 400, 850},
		{"price too large", "1", `{"price":2147483648}`, 400, 850},
		{"missing price", "1", `{}`, 400, 850},
		{"null price", "1", `{"price":null}`, 400, 850},
		{"decimal price", "1", `{"price":950.5}`, 400, 850},
		{"string price", "1", `{"price":"950"}`, 400, 850},
		{"unknown field", "1", `{"price":950,"name":"Changed"}`, 400, 850},
		{"multiple values", "1", `{"price":950} {}`, 400, 850},
		{"empty body", "1", ``, 400, 850},
		{"invalid ID", "abc", `{"price":950}`, 400, 850},
		{"zero ID", "0", `{"price":950}`, 400, 850},
		{"missing dish", "999", `{"price":950}`, 404, 850},
		{
			"oversized body",
			"1",
			`{"price":950}` + strings.Repeat(" ", 64*1024),
			413,
			850,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dishes := repository.NewDishRepository()
			restaurants := repository.NewRestaurantRepository()
			handler := NewDishHandler(dishes, restaurants)

			before, found, err := dishes.FindByID(context.Background(), 1)
			if err != nil || !found {
				t.Fatalf("read original dish: found=%v, err=%v", found, err)
			}

			router := chi.NewRouter()
			router.Patch("/api/dishes/{id}/price", handler.UpdatePrice)

			request := httptest.NewRequest(
				http.MethodPatch,
				"/api/dishes/"+tt.id+"/price",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")
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

			want := before
			want.Price = tt.wantPrice

			saved, found, err := dishes.FindByID(context.Background(), 1)
			if err != nil || !found {
				t.Fatalf("read saved dish: found=%v, err=%v", found, err)
			}
			if saved != want {
				t.Errorf("expected saved dish %+v, got %+v", want, saved)
			}

			if tt.wantStatus == http.StatusOK {
				var result models.Dish
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if result != want {
					t.Errorf("expected response %+v, got %+v", want, result)
				}
			}

			other, found, err := dishes.FindByID(context.Background(), 2)
			if err != nil || !found {
				t.Fatalf("read other dish: found=%v, err=%v", found, err)
			}
			if other.Price != 650 {
				t.Errorf("another dish's price changed: got %d", other.Price)
			}
		})
	}
}

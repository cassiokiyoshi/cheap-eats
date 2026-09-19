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

func TestCreateRestaurantBilingualNames(t *testing.T) {
	tests := []struct {
		name       string
		extra      string
		wantStatus int
		wantJA     string
		wantEN     string
	}{
		{
			name:       "both names are trimmed",
			extra:      `,"name_ja":"  さくら食堂  ","name_en":"  Sakura Shokudo  "`,
			wantStatus: http.StatusCreated,
			wantJA:     "さくら食堂",
			wantEN:     "Sakura Shokudo",
		},
		{
			name:       "names omitted",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "explicit null names",
			extra:      `,"name_ja":null,"name_en":null`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Japanese name only",
			extra:      `,"name_ja":"さくら食堂"`,
			wantStatus: http.StatusCreated,
			wantJA:     "さくら食堂",
		},
		{
			name:       "English name only",
			extra:      `,"name_en":"Sakura Shokudo"`,
			wantStatus: http.StatusCreated,
			wantEN:     "Sakura Shokudo",
		},
		{
			name:       "blank Japanese name",
			extra:      `,"name_ja":"　 "`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty English name",
			extra:      `,"name_en":""`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong name type",
			extra:      `,"name_en":123`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := repository.NewRestaurantRepository()
			handler := NewRestaurantHandler(store)

			before, err := store.List(ctx)
			if err != nil {
				t.Fatal(err)
			}
			beforeCount := len(before)

			body := `{
				"name":"Sakura Shokudo",
				"address":"Tokyo",
				"latitude":35.6812,
				"longitude":139.7671` + tt.extra + `}`

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/restaurants",
				strings.NewReader(body),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler.Create(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf(
					"expected %d, got %d: %s",
					tt.wantStatus,
					response.Code,
					response.Body.String(),
				)
			}

			after, err := store.List(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantStatus == http.StatusBadRequest {
				if len(after) != beforeCount {
					t.Fatal("invalid input saved a restaurant")
				}
				return
			}

			if len(after) != beforeCount+1 {
				t.Fatal("expected exactly one new restaurant")
			}

			var created models.Restaurant
			if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			saved, found, err := store.FindByID(ctx, created.ID)
			if err != nil || !found {
				t.Fatalf("find restaurant: found=%v, err=%v", found, err)
			}

			for label, restaurant := range map[string]models.Restaurant{
				"response": created,
				"stored":   saved,
			} {
				if restaurant.Name != "Sakura Shokudo" {
					t.Errorf("%s: original name changed", label)
				}

				assertBilingualName(
					t, label+" Japanese", restaurant.NameJA, tt.wantJA,
				)
				assertBilingualName(
					t, label+" English", restaurant.NameEN, tt.wantEN,
				)
			}
		})
	}
}

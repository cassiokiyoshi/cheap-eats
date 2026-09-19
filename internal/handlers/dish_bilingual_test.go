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

func TestCreateDishBilingualNames(t *testing.T) {
	tests := []struct {
		name       string
		extra      string
		wantStatus int
		wantJA     string
		wantEN     string
	}{
		{
			name:       "both names are trimmed",
			extra:      `,"name_ja":"  醤油ラーメン  ","name_en":"  Shoyu Ramen  "`,
			wantStatus: http.StatusCreated,
			wantJA:     "醤油ラーメン",
			wantEN:     "Shoyu Ramen",
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
			extra:      `,"name_ja":"醤油ラーメン"`,
			wantStatus: http.StatusCreated,
			wantJA:     "醤油ラーメン",
		},
		{
			name:       "English name only",
			extra:      `,"name_en":"Shoyu Ramen"`,
			wantStatus: http.StatusCreated,
			wantEN:     "Shoyu Ramen",
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
			dishes := repository.NewDishRepository()
			handler := NewDishHandler(
				dishes,
				repository.NewRestaurantRepository(),
			)

			before, err := dishes.List(ctx)
			if err != nil {
				t.Fatal(err)
			}
			beforeCount := len(before)

			body := `{"restaurant_id":1,"name":"Shoyu Ramen","price":850` +
				tt.extra + `}`

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/dishes",
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

			after, err := dishes.List(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantStatus == http.StatusBadRequest {
				if len(after) != beforeCount {
					t.Fatal("invalid input saved a dish")
				}
				return
			}

			if len(after) != beforeCount+1 {
				t.Fatal("expected exactly one new dish")
			}

			var created models.Dish
			if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			saved, found, err := dishes.FindByID(ctx, created.ID)
			if err != nil || !found {
				t.Fatalf("find dish: found=%v, err=%v", found, err)
			}

			for label, dish := range map[string]models.Dish{
				"response": created,
				"stored":   saved,
			} {
				if dish.Name != "Shoyu Ramen" {
					t.Errorf("%s: original name changed", label)
				}

				assertBilingualName(t, label+" Japanese", dish.NameJA, tt.wantJA)
				assertBilingualName(t, label+" English", dish.NameEN, tt.wantEN)
			}
		})
	}
}

func assertBilingualName(t *testing.T, label string, got *string, want string) {
	t.Helper()

	if want == "" {
		if got != nil {
			t.Errorf("%s: expected nil, got %q", label, *got)
		}
		return
	}

	if got == nil {
		t.Errorf("%s: expected %q, got nil", label, want)
		return
	}

	if *got != want {
		t.Errorf("%s: expected %q, got %q", label, want, *got)
	}
}

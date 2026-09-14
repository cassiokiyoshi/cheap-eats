package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNearbyRejectsInvalidCoordinates(t *testing.T) {
	dishHandler := NewDishSearchHandler(nil)
	restaurantHandler := NewRestaurantHandler(nil)

	handlers := []struct {
		name string
		path string
		call http.HandlerFunc
	}{
		{
			name: "dishes",
			path: "/api/dishes/nearby",
			call: dishHandler.Nearby,
		},
		{
			name: "restaurants",
			path: "/api/restaurants/nearby",
			call: restaurantHandler.Nearby,
		},
	}

	cases := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "NaN latitude",
			query: "lat=NaN&lng=139.7671&max_price=1000",
			want:  "invalid latitude",
		},
		{
			name:  "infinite longitude",
			query: "lat=35.6812&lng=Inf&max_price=1000",
			want:  "invalid longitude",
		},
		{
			name:  "radius below minimum",
			query: "lat=35.6812&lng=139.7671&radius=0.5&max_price=1000",
			want:  "radius must be between 1 and 5000 meters",
		},
		{
			name:  "NaN radius",
			query: "lat=35.6812&lng=139.7671&radius=NaN&max_price=1000",
			want:  "radius must be between 1 and 5000 meters",
		},
	}

	for _, handler := range handlers {
		for _, tc := range cases {
			t.Run(handler.name+"/"+tc.name, func(t *testing.T) {
				request := httptest.NewRequest(
					http.MethodGet,
					handler.path+"?"+tc.query,
					nil,
				)
				response := httptest.NewRecorder()

				handler.call(response, request)

				if response.Code != http.StatusBadRequest {
					t.Fatalf(
						"expected status 400, got %d",
						response.Code,
					)
				}

				if got := strings.TrimSpace(response.Body.String()); got != tc.want {
					t.Errorf(
						"expected message %q, got %q",
						tc.want,
						got,
					)
				}
			})
		}
	}
}

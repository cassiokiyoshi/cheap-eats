package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerErrorHidesInternalDetails(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/dishes",
		nil,
	)
	response := httptest.NewRecorder()

	internalError := errors.New(
		"database connection failed: private internal details",
	)

	serverError(response, request, internalError)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			response.Code,
		)
	}

	wantBody := "internal server error\n"
	if response.Body.String() != wantBody {
		t.Errorf(
			"expected body %q, got %q",
			wantBody,
			response.Body.String(),
		)
	}
}

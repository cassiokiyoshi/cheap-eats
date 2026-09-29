package pathfinder

import (
	"strings"
	"testing"
)

func TestSearchRequestValidate(t *testing.T) {
	validRequest := SearchRequest{
		Algorithm: AlgorithmAStar,
		Rows:      10,
		Cols:      12,
		Start:     Position{Row: 1, Col: 1},
		Goal:      Position{Row: 8, Col: 10},
		Walls: []Position{
			{Row: 3, Col: 4},
		},
	}

	tests := []struct {
		name            string
		request         SearchRequest
		expectedMessage string
	}{
		{
			name:    "valid request",
			request: validRequest,
		},
		{
			name: "unsupported algorithm",
			request: func() SearchRequest {
				request := validRequest
				request.Algorithm = "unknown"
				return request
			}(),
			expectedMessage: "unsupported algorithm",
		},
		{
			name: "grid is too small",
			request: func() SearchRequest {
				request := validRequest
				request.Rows = 1
				return request
			}(),
			expectedMessage: "rows must be between",
		},
		{
			name: "start is outside grid",
			request: func() SearchRequest {
				request := validRequest
				request.Start = Position{Row: -1, Col: 0}
				return request
			}(),
			expectedMessage: "start position is outside",
		},
		{
			name: "goal is outside grid",
			request: func() SearchRequest {
				request := validRequest
				request.Goal = Position{Row: 20, Col: 0}
				return request
			}(),
			expectedMessage: "goal position is outside",
		},
		{
			name: "start equals goal",
			request: func() SearchRequest {
				request := validRequest
				request.Goal = request.Start
				return request
			}(),
			expectedMessage: "must be different",
		},
		{
			name: "wall covers start",
			request: func() SearchRequest {
				request := validRequest
				request.Walls = []Position{request.Start}
				return request
			}(),
			expectedMessage: "start position cannot be a wall",
		},
		{
			name: "wall is outside grid",
			request: func() SearchRequest {
				request := validRequest
				request.Walls = []Position{{Row: 100, Col: 100}}
				return request
			}(),
			expectedMessage: "wall at row 100",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.request.Validate()

			if test.expectedMessage == "" {
				if err != nil {
					t.Fatalf("expected request to be valid, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q", test.expectedMessage)
			}

			if !strings.Contains(err.Error(), test.expectedMessage) {
				t.Fatalf(
					"expected error containing %q, got %q",
					test.expectedMessage,
					err.Error(),
				)
			}
		})
	}
}

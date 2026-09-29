package pathfinder

import (
	"context"
	"testing"
)

func TestAlgorithmsFindSameShortestPathLength(t *testing.T) {
	algorithms := []Algorithm{
		AlgorithmBFS,
		AlgorithmDijkstra,
		AlgorithmAStar,
	}

	for _, algorithm := range algorithms {
		t.Run(string(algorithm), func(t *testing.T) {
			request := SearchRequest{
				Algorithm: algorithm,
				Rows:      7,
				Cols:      7,
				Start:     Position{Row: 3, Col: 0},
				Goal:      Position{Row: 3, Col: 6},
				Walls: []Position{
					{Row: 1, Col: 3},
					{Row: 2, Col: 3},
					{Row: 3, Col: 3},
					{Row: 4, Col: 3},
					{Row: 5, Col: 3},
					{Row: 6, Col: 3},
				},
			}

			path, found, err := Run(
				context.Background(),
				request,
				nil,
			)

			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if !found {
				t.Fatal("expected a path to be found")
			}

			const expectedPathLength = 13

			if len(path) != expectedPathLength {
				t.Fatalf(
					"expected path length %d, got %d",
					expectedPathLength,
					len(path),
				)
			}

			if path[0] != request.Start {
				t.Fatalf(
					"path starts at %v instead of %v",
					path[0],
					request.Start,
				)
			}

			if path[len(path)-1] != request.Goal {
				t.Fatalf(
					"path ends at %v instead of %v",
					path[len(path)-1],
					request.Goal,
				)
			}
		})
	}
}

func TestRunRejectsInvalidRequest(t *testing.T) {
	request := SearchRequest{
		Algorithm: AlgorithmBFS,
		Rows:      1,
		Cols:      10,
		Start:     Position{Row: 0, Col: 0},
		Goal:      Position{Row: 0, Col: 9},
	}

	_, _, err := Run(context.Background(), request, nil)

	if err == nil {
		t.Fatal("expected invalid request to return an error")
	}
}

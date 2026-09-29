package pathfinder

import (
	"context"
	"errors"
	"testing"
)

func TestRunBFSFindsShortestPath(t *testing.T) {
	request := SearchRequest{
		Algorithm: AlgorithmBFS,
		Rows:      3,
		Cols:      5,
		Start:     Position{Row: 1, Col: 0},
		Goal:      Position{Row: 1, Col: 4},
	}

	var visited []Position

	path, found, err := RunBFS(
		context.Background(),
		request,
		func(position Position) error {
			visited = append(visited, position)
			return nil
		},
	)

	if err != nil {
		t.Fatalf("RunBFS returned an error: %v", err)
	}

	if !found {
		t.Fatal("RunBFS did not find a path")
	}

	if len(path) != 5 {
		t.Fatalf("expected path length 5, got %d", len(path))
	}

	if path[0] != request.Start {
		t.Fatalf("path starts at %v instead of %v", path[0], request.Start)
	}

	if path[len(path)-1] != request.Goal {
		t.Fatalf("path ends at %v instead of %v", path[len(path)-1], request.Goal)
	}

	if len(visited) == 0 {
		t.Fatal("expected at least one visited position")
	}
}

func TestRunBFSAvoidsWalls(t *testing.T) {
	wall := Position{Row: 1, Col: 1}

	request := SearchRequest{
		Algorithm: AlgorithmBFS,
		Rows:      3,
		Cols:      4,
		Start:     Position{Row: 1, Col: 0},
		Goal:      Position{Row: 1, Col: 3},
		Walls:     []Position{wall},
	}

	path, found, err := RunBFS(context.Background(), request, nil)
	if err != nil {
		t.Fatalf("RunBFS returned an error: %v", err)
	}

	if !found {
		t.Fatal("RunBFS did not find a route around the wall")
	}

	for _, position := range path {
		if position == wall {
			t.Fatalf("path passes through wall at %v", wall)
		}
	}

	if len(path) != 6 {
		t.Fatalf("expected path length 6, got %d", len(path))
	}
}

func TestRunBFSReportsNoPath(t *testing.T) {
	request := SearchRequest{
		Algorithm: AlgorithmBFS,
		Rows:      3,
		Cols:      3,
		Start:     Position{Row: 0, Col: 0},
		Goal:      Position{Row: 2, Col: 2},
		Walls: []Position{
			{Row: 1, Col: 0},
			{Row: 1, Col: 1},
			{Row: 1, Col: 2},
		},
	}

	path, found, err := RunBFS(context.Background(), request, nil)
	if err != nil {
		t.Fatalf("RunBFS returned an error: %v", err)
	}

	if found {
		t.Fatal("RunBFS unexpectedly found a path")
	}

	if path != nil {
		t.Fatalf("expected no path, got %v", path)
	}
}

func TestRunBFSHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	request := SearchRequest{
		Algorithm: AlgorithmBFS,
		Rows:      10,
		Cols:      10,
		Start:     Position{Row: 0, Col: 0},
		Goal:      Position{Row: 9, Col: 9},
	}

	_, _, err := RunBFS(ctx, request, nil)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

package pathfinder

import "fmt"

const (
	minGridSize = 2
	maxGridSize = 200
)

func (request SearchRequest) Validate() error {
	if !request.Algorithm.IsValid() {
		return fmt.Errorf("unsupported algorithm %q", request.Algorithm)
	}

	if request.Rows < minGridSize || request.Rows > maxGridSize {
		return fmt.Errorf(
			"rows must be between %d and %d",
			minGridSize,
			maxGridSize,
		)
	}

	if request.Cols < minGridSize || request.Cols > maxGridSize {
		return fmt.Errorf(
			"columns must be between %d and %d",
			minGridSize,
			maxGridSize,
		)
	}

	if !request.contains(request.Start) {
		return fmt.Errorf("start position is outside the grid")
	}

	if !request.contains(request.Goal) {
		return fmt.Errorf("goal position is outside the grid")
	}

	if request.Start == request.Goal {
		return fmt.Errorf("start and goal positions must be different")
	}

	for _, wall := range request.Walls {
		if !request.contains(wall) {
			return fmt.Errorf(
				"wall at row %d, column %d is outside the grid",
				wall.Row,
				wall.Col,
			)
		}

		if wall == request.Start {
			return fmt.Errorf("start position cannot be a wall")
		}

		if wall == request.Goal {
			return fmt.Errorf("goal position cannot be a wall")
		}
	}

	return nil
}

func (algorithm Algorithm) IsValid() bool {
	switch algorithm {
	case AlgorithmAStar, AlgorithmDijkstra, AlgorithmBFS:
		return true
	default:
		return false
	}
}

func (request SearchRequest) contains(position Position) bool {
	return position.Row >= 0 &&
		position.Row < request.Rows &&
		position.Col >= 0 &&
		position.Col < request.Cols
}

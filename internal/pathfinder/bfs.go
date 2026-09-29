package pathfinder

import "context"

type VisitFunc func(Position) error

func RunBFS(
	ctx context.Context,
	request SearchRequest,
	onVisit VisitFunc,
) ([]Position, bool, error) {
	walls := make(map[Position]struct{}, len(request.Walls))
	for _, wall := range request.Walls {
		walls[wall] = struct{}{}
	}

	queue := []Position{request.Start}
	visited := map[Position]bool{
		request.Start: true,
	}
	previous := make(map[Position]Position)

	for len(queue) > 0 {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		current := queue[0]
		queue = queue[1:]

		if onVisit != nil {
			if err := onVisit(current); err != nil {
				return nil, false, err
			}
		}

		if current == request.Goal {
			return reconstructPath(previous, request.Start, request.Goal), true, nil
		}

		for _, neighbor := range neighbors(current) {
			if !request.contains(neighbor) {
				continue
			}

			if visited[neighbor] {
				continue
			}

			if _, blocked := walls[neighbor]; blocked {
				continue
			}

			visited[neighbor] = true
			previous[neighbor] = current
			queue = append(queue, neighbor)
		}
	}

	return nil, false, nil
}

func neighbors(position Position) []Position {
	return []Position{
		{Row: position.Row - 1, Col: position.Col},
		{Row: position.Row, Col: position.Col + 1},
		{Row: position.Row + 1, Col: position.Col},
		{Row: position.Row, Col: position.Col - 1},
	}
}

func reconstructPath(
	previous map[Position]Position,
	start Position,
	goal Position,
) []Position {
	path := []Position{goal}
	current := goal

	for current != start {
		current = previous[current]
		path = append(path, current)
	}

	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}

	return path
}

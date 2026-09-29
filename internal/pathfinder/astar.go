package pathfinder

import (
	"container/heap"
	"context"
)

func RunAStar(
	ctx context.Context,
	request SearchRequest,
	onVisit VisitFunc,
) ([]Position, bool, error) {
	walls := make(map[Position]struct{}, len(request.Walls))
	for _, wall := range request.Walls {
		walls[wall] = struct{}{}
	}

	distances := map[Position]int{
		request.Start: 0,
	}
	previous := make(map[Position]Position)
	processed := make(map[Position]bool)

	queue := &priorityQueue{}
	heap.Init(queue)

	nextOrder := 0
	heap.Push(queue, queueItem{
		position: request.Start,
		priority: manhattanDistance(request.Start, request.Goal),
		order:    nextOrder,
	})

	for queue.Len() > 0 {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		item := heap.Pop(queue).(queueItem)
		current := item.position

		if processed[current] {
			continue
		}
		processed[current] = true

		if onVisit != nil {
			if err := onVisit(current); err != nil {
				return nil, false, err
			}
		}

		if current == request.Goal {
			return reconstructPath(previous, request.Start, request.Goal), true, nil
		}

		currentDistance := distances[current]

		for _, neighbor := range neighbors(current) {
			if !request.contains(neighbor) {
				continue
			}

			if _, blocked := walls[neighbor]; blocked {
				continue
			}

			newDistance := currentDistance + 1
			oldDistance, discovered := distances[neighbor]

			if discovered && newDistance >= oldDistance {
				continue
			}

			distances[neighbor] = newDistance
			previous[neighbor] = current
			nextOrder++

			estimatedTotalDistance :=
				newDistance + manhattanDistance(neighbor, request.Goal)

			heap.Push(queue, queueItem{
				position: neighbor,
				priority: estimatedTotalDistance,
				order:    nextOrder,
			})
		}
	}

	return nil, false, nil
}

func manhattanDistance(from Position, to Position) int {
	rowDistance := absolute(from.Row - to.Row)
	colDistance := absolute(from.Col - to.Col)
	return rowDistance + colDistance
}

func absolute(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

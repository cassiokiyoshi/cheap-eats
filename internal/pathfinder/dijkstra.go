package pathfinder

import (
	"container/heap"
	"context"
)

type queueItem struct {
	position Position
	priority int
	order    int
}

type priorityQueue []queueItem

func (queue priorityQueue) Len() int {
	return len(queue)
}

func (queue priorityQueue) Less(i int, j int) bool {
	if queue[i].priority == queue[j].priority {
		return queue[i].order < queue[j].order
	}

	return queue[i].priority < queue[j].priority
}

func (queue priorityQueue) Swap(i int, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue *priorityQueue) Push(value any) {
	*queue = append(*queue, value.(queueItem))
}

func (queue *priorityQueue) Pop() any {
	old := *queue
	lastIndex := len(old) - 1
	item := old[lastIndex]
	*queue = old[:lastIndex]
	return item
}

func RunDijkstra(
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
		priority: 0,
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

			heap.Push(queue, queueItem{
				position: neighbor,
				priority: newDistance,
				order:    nextOrder,
			})
		}
	}

	return nil, false, nil
}

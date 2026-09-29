package pathfinder

import (
	"context"
	"fmt"
)

func Run(
	ctx context.Context,
	request SearchRequest,
	onVisit VisitFunc,
) ([]Position, bool, error) {
	if err := request.Validate(); err != nil {
		return nil, false, err
	}

	switch request.Algorithm {
	case AlgorithmBFS:
		return RunBFS(ctx, request, onVisit)

	case AlgorithmDijkstra:
		return RunDijkstra(ctx, request, onVisit)

	case AlgorithmAStar:
		return RunAStar(ctx, request, onVisit)

	default:
		return nil, false, fmt.Errorf(
			"unsupported algorithm %q",
			request.Algorithm,
		)
	}
}

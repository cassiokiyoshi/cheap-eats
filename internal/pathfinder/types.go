package pathfinder

type Algorithm string

const (
	AlgorithmAStar    Algorithm = "astar"
	AlgorithmDijkstra Algorithm = "dijkstra"
	AlgorithmBFS      Algorithm = "bfs"
)

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type SearchRequest struct {
	Algorithm Algorithm  `json:"algorithm"`
	Rows      int        `json:"rows"`
	Cols      int        `json:"cols"`
	Start     Position   `json:"start"`
	Goal      Position   `json:"goal"`
	Walls     []Position `json:"walls"`
}

type StreamEvent struct {
	Type       string     `json:"type"`
	Position   *Position  `json:"position,omitempty"`
	Positions  []Position `json:"positions,omitempty"`
	Visited    int        `json:"visited,omitempty"`
	PathLength int        `json:"pathLength,omitempty"`
	DurationMS float64    `json:"durationMs,omitempty"`
	Message    string     `json:"message,omitempty"`
	Found      *bool      `json:"found,omitempty"`
}

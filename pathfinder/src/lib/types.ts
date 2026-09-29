export type Algorithm = 'astar' | 'dijkstra' | 'bfs';

export type Position = {
	row: number;
	col: number;
};

export type SearchRequest = {
	algorithm: Algorithm;
	rows: number;
	cols: number;
	start: Position;
	goal: Position;
	walls: Position[];
};

export type StreamEvent =
	| {
			type: 'visited';
			position: Position;
	  }
	| {
			type: 'path';
			positions: Position[];
	  }
	| {
			type: 'complete';
			visited: number;
			pathLength: number;
			durationMs: number;
			found: boolean;
	  }
	| {
			type: 'error';
			message: string;
	  };

export type SearchStats = {
	visited: number;
	pathLength: number;
	durationMs: number;
	found: boolean;
};

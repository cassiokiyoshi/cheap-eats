import type { Position } from '$lib/types';

export const GRID_ROWS = 20;
export const GRID_COLS = 32;

export const DEFAULT_START: Position = {
	row: Math.floor(GRID_ROWS / 2),
	col: 5
};

export const DEFAULT_GOAL: Position = {
	row: Math.floor(GRID_ROWS / 2),
	col: GRID_COLS - 6
};

export function positionKey(position: Position): string {
	return `${position.row}:${position.col}`;
}

export function positionsAreEqual(
	first: Position,
	second: Position
): boolean {
	return first.row === second.row && first.col === second.col;
}

export function createGridPositions(
	rows: number,
	cols: number
): Position[] {
	const positions: Position[] = [];

	for (let row = 0; row < rows; row++) {
		for (let col = 0; col < cols; col++) {
			positions.push({ row, col });
		}
	}

	return positions;
}

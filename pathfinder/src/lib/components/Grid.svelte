<script lang="ts">
	import {
		createGridPositions,
		positionKey,
		positionsAreEqual
	} from '$lib/grid';
	import type { Position } from '$lib/types';

	type Endpoint = 'start' | 'goal';

	type GridProps = {
		rows: number;
		cols: number;
		start: Position;
		goal: Position;
		walls: Set<string>;
		explored: Set<string>;
		path: Set<string>;
		disabled: boolean;
		onWallChange: (position: Position, isWall: boolean) => void;
		onEndpointChange: (
			endpoint: Endpoint,
			position: Position
		) => void;
	};

	let {
		rows,
		cols,
		start,
		goal,
		walls,
		explored,
		path,
		disabled,
		onWallChange,
		onEndpointChange
	}: GridProps = $props();

	let gridElement: HTMLDivElement;
	let painting = false;
	let paintAsWall = true;
	let draggedEndpoint: Endpoint | null = null;
	let lastPaintedKey = '';

	const positions = $derived(createGridPositions(rows, cols));

	function cellState(position: Position): string {
		if (positionsAreEqual(position, start)) {
			return 'start';
		}

		if (positionsAreEqual(position, goal)) {
			return 'goal';
		}

		const key = positionKey(position);

		if (path.has(key)) {
			return 'path';
		}

		if (walls.has(key)) {
			return 'wall';
		}

		if (explored.has(key)) {
			return 'explored';
		}

		return 'empty';
	}

	function cellLabel(position: Position): string {
		return `Row ${position.row + 1}, column ${position.col + 1}: ${cellState(position)}`;
	}

	function startInteraction(
		event: PointerEvent,
		position: Position
	): void {
		if (disabled) {
			return;
		}

		event.preventDefault();
		lastPaintedKey = '';

		if (positionsAreEqual(position, start)) {
			draggedEndpoint = 'start';
			return;
		}

		if (positionsAreEqual(position, goal)) {
			draggedEndpoint = 'goal';
			return;
		}

		painting = true;
		paintAsWall = !walls.has(positionKey(position));
		paint(position);
	}

	function continueInteraction(event: PointerEvent): void {
		if (disabled || (!painting && !draggedEndpoint)) {
			return;
		}

		const element = document.elementFromPoint(
			event.clientX,
			event.clientY
		);

		const cell = element?.closest<HTMLButtonElement>('[data-grid-cell]');

		if (!cell || !gridElement.contains(cell)) {
			return;
		}

		const row = Number(cell.dataset.row);
		const col = Number(cell.dataset.col);

		if (!Number.isInteger(row) || !Number.isInteger(col)) {
			return;
		}

		const position = { row, col };

		if (draggedEndpoint) {
			moveEndpoint(draggedEndpoint, position);
			return;
		}

		paint(position);
	}

	function moveEndpoint(
		endpoint: Endpoint,
		position: Position
	): void {
		const key = positionKey(position);
		const otherEndpoint = endpoint === 'start' ? goal : start;

		if (
			key === lastPaintedKey ||
			walls.has(key) ||
			positionsAreEqual(position, otherEndpoint)
		) {
			return;
		}

		lastPaintedKey = key;
		onEndpointChange(endpoint, position);
	}

	function paint(position: Position): void {
		if (
			positionsAreEqual(position, start) ||
			positionsAreEqual(position, goal)
		) {
			return;
		}

		const key = positionKey(position);

		if (key === lastPaintedKey) {
			return;
		}

		lastPaintedKey = key;
		onWallChange(position, paintAsWall);
	}

	function stopInteraction(): void {
		painting = false;
		draggedEndpoint = null;
		lastPaintedKey = '';
	}

	function handleKeydown(
		event: KeyboardEvent,
		position: Position
	): void {
		if (
			disabled ||
			positionsAreEqual(position, start) ||
			positionsAreEqual(position, goal)
		) {
			return;
		}

		if (event.key !== 'Enter' && event.key !== ' ') {
			return;
		}

		event.preventDefault();
		onWallChange(position, !walls.has(positionKey(position)));
	}
</script>

<svelte:window
	onpointerup={stopInteraction}
	onpointercancel={stopInteraction}
/>

<div
	bind:this={gridElement}
	class="grid"
	class:disabled
	style:--columns={cols}
	role="grid"
  tabindex="0"
	aria-label="Pathfinding grid"
	onpointermove={continueInteraction}
>
	{#each positions as position (`${position.row}:${position.col}`)}
		<button
			class="cell"
			class:start={cellState(position) === 'start'}
			class:goal={cellState(position) === 'goal'}
			class:wall={cellState(position) === 'wall'}
			class:explored={cellState(position) === 'explored'}
			class:path={cellState(position) === 'path'}
			type="button"
			role="gridcell"
			data-grid-cell
			data-row={position.row}
			data-col={position.col}
			aria-label={cellLabel(position)}
			disabled={disabled}
			onpointerdown={(event) => startInteraction(event, position)}
			onkeydown={(event) => handleKeydown(event, position)}
		>
			{#if cellState(position) === 'start'}
				<span aria-hidden="true">S</span>
			{:else if cellState(position) === 'goal'}
				<span aria-hidden="true">G</span>
			{/if}
		</button>
	{/each}
</div>

<style>
	.grid {
		--grid-line: #26334a;
		display: grid;
		grid-template-columns: repeat(var(--columns), minmax(0, 1fr));
		width: 100%;
		overflow: hidden;
		border: 1px solid var(--grid-line);
		border-radius: 0.75rem;
		background: var(--grid-line);
		gap: 1px;
		touch-action: none;
		user-select: none;
	}

	.grid.disabled {
		opacity: 0.82;
	}

	.cell {
		display: grid;
		min-width: 0;
		aspect-ratio: 1;
		place-items: center;
		padding: 0;
		border: 0;
		background: #101827;
		color: white;
		font: inherit;
		font-size: clamp(0.45rem, 1.4vw, 0.75rem);
		cursor: crosshair;
	}

	.cell:hover {
		background: #19243a;
	}

	.cell:focus-visible {
		position: relative;
		z-index: 1;
		outline: 2px solid white;
		outline-offset: -2px;
	}

	.cell.wall {
		background: #546178;
	}

	.cell.explored {
		background: #185b78;
		animation: appear 180ms ease-out;
	}

	.cell.path {
		background: #f3c745;
		animation: appear 180ms ease-out;
	}

	.cell.start {
		background: #24a76a;
		font-weight: 800;
		cursor: grab;
	}

	.cell.goal {
		background: #e55353;
		font-weight: 800;
		cursor: grab;
	}

	.cell:disabled {
		cursor: not-allowed;
	}

	@keyframes appear {
		from {
			opacity: 0.25;
			transform: scale(0.65);
		}

		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.cell {
			animation: none;
		}
	}
</style>

<script lang="ts">
	import { onDestroy } from 'svelte';
	import Controls from '$lib/components/Controls.svelte';
  import Stats from '$lib/components/Stats.svelte';
	import Grid from '$lib/components/Grid.svelte';
	import {
		DEFAULT_GOAL,
		DEFAULT_START,
		GRID_COLS,
		GRID_ROWS,
		positionKey
	} from '$lib/grid';
	import { streamSearch } from '$lib/search';
	import type {
    Algorithm,
    Position,
    SearchRequest,
    SearchStats,
    StreamEvent
  } from '$lib/types';

	let start = $state({ ...DEFAULT_START });
	let goal = $state({ ...DEFAULT_GOAL });
	let walls = $state(new Set<string>());
	let explored = $state(new Set<string>());
	let path = $state(new Set<string>());

	let algorithm = $state<Algorithm>('astar');
	let speed = $state(20);
	let isRunning = $state(false);
	let isPaused = $state(false);
	let status = $state('Draw walls, then run a search.');

	let abortController: AbortController | null = null;

  let stats = $state<SearchStats | null>(null);

	onDestroy(() => {
		abortController?.abort();
	});

	function updateWall(position: Position, isWall: boolean): void {
		const nextWalls = new Set(walls);
		const key = positionKey(position);

		if (isWall) {
			nextWalls.add(key);
		} else {
			nextWalls.delete(key);
		}

		walls = nextWalls;
		clearVisualization();
	}

	function clearVisualization(): void {
		explored = new Set();
		path = new Set();
    stats = null;
	}

	function clearPath(): void {
		clearVisualization();
		status = 'Visualization cleared.';
	}

	function clearWalls(): void {
		walls = new Set();
		clearVisualization();
		status = 'Walls cleared.';
	}

	function resetGrid(): void {
		abortController?.abort();

		start = { ...DEFAULT_START };
		goal = { ...DEFAULT_GOAL };
		walls = new Set();
		clearVisualization();

		algorithm = 'astar';
		speed = 20;
		isRunning = false;
		isPaused = false;
		status = 'Grid reset.';
	}

	function togglePause(): void {
		isPaused = !isPaused;
		status = isPaused ? 'Visualization paused.' : 'Visualization resumed.';
	}

  function stopSearch(): void {
    if (!abortController) {
      return;
    }

    abortController.abort();
    isPaused = false;
    status = 'Search stopped.';
  }

	async function startSearch(): Promise<void> {
		clearVisualization();

		isRunning = true;
		isPaused = false;
		status = `Running ${algorithmLabel(algorithm)}…`;

		const controller = new AbortController();
		abortController = controller;

		const eventQueue: StreamEvent[] = [];
		let streamFinished = false;
		let streamError: unknown = null;

		const request: SearchRequest = {
			algorithm,
			rows: GRID_ROWS,
			cols: GRID_COLS,
			start,
			goal,
			walls: Array.from(walls, positionFromKey)
		};

		const receiveEvents = streamSearch(
			request,
			controller.signal,
			(event) => {
				eventQueue.push(event);
			}
		)
			.catch((error: unknown) => {
				streamError = error;
			})
			.finally(() => {
				streamFinished = true;
			});

		try {
			await playEvents(
				eventQueue,
				() => streamFinished,
				controller.signal
			);

			await receiveEvents;

			if (streamError) {
				throw streamError;
			}
		} catch (error: unknown) {
			if (!controller.signal.aborted) {
				status =
					error instanceof Error
						? error.message
						: 'The search could not be completed.';
			}
		} finally {
			if (abortController === controller) {
				abortController = null;
				isRunning = false;
				isPaused = false;
			}
		}
	}

	async function playEvents(
		queue: StreamEvent[],
		hasFinished: () => boolean,
		signal: AbortSignal
	): Promise<void> {
		while (!hasFinished() || queue.length > 0) {
			if (signal.aborted) {
				return;
			}

			if (isPaused) {
				await delay(40);
				continue;
			}

			const event = queue.shift();

			if (!event) {
				await delay(5);
				continue;
			}

			if (event.type === 'path') {
        await animatePath(event.positions, signal);
        continue;
      }

      applyEvent(event);

      if (event.type === 'visited') {
        await delay(speed);
      }
    }
  }

  async function animatePath(
    positions: Position[],
    signal: AbortSignal
  ): Promise<void> {
    for (const position of positions) {
      while (isPaused && !signal.aborted) {
        await delay(40);
      }

      if (signal.aborted) {
        return;
      }

      path = new Set(path).add(positionKey(position));
      await delay(Math.max(20, speed));
    }
  }

	function applyEvent(event: StreamEvent): void {
		switch (event.type) {
			case 'visited':
				explored = new Set(explored).add(
					positionKey(event.position)
				);
				break;

			case 'complete':
        stats = {
          visited: event.visited,
          pathLength: event.pathLength,
          durationMs: event.durationMs,
          found: event.found
        };

        status = event.found
          ? `Path found in ${event.pathLength} steps.`
          : 'No path exists between the start and goal.';
        break;

			case 'error':
				status = event.message;
				break;
		}
	}

	function positionFromKey(key: string): Position {
		const [row, col] = key.split(':').map(Number);
		return { row, col };
	}

	function algorithmLabel(value: Algorithm): string {
		switch (value) {
			case 'astar':
				return 'A* search';
			case 'dijkstra':
				return 'Dijkstra';
			case 'bfs':
				return 'breadth-first search';
		}
	}

	function delay(milliseconds: number): Promise<void> {
		return new Promise((resolve) => {
			window.setTimeout(resolve, milliseconds);
		});
	}

  type Endpoint = 'start' | 'goal';

  function updateEndpoint(
    endpoint: Endpoint,
    position: Position
  ): void {
    if (endpoint === 'start') {
      start = { ...position };
    } else {
      goal = { ...position };
    }

    clearVisualization();
    status = `${endpoint === 'start' ? 'Start' : 'Goal'} moved.`;
  }
</script>

<svelte:head>
	<title>Pathfinder</title>
	<meta
		name="description"
		content="Explore A*, Dijkstra, and breadth-first search on an editable grid."
	/>
</svelte:head>

<div class="page-shell">
	<header class="hero">
		<div>
			<p class="eyebrow">Interactive algorithm laboratory</p>
			<h1>Pathfinder</h1>
			<p class="introduction">
				Draw obstacles, choose an algorithm, and watch the search unfold
				cell by cell.
			</p>
		</div>

		<div class="legend" aria-label="Grid legend">
			<span><i class="start-swatch"></i> Start</span>
			<span><i class="goal-swatch"></i> Goal</span>
			<span><i class="wall-swatch"></i> Wall</span>
			<span><i class="visited-swatch"></i> Explored</span>
			<span><i class="path-swatch"></i> Path</span>
		</div>
	</header>

	<main>
		<section class="workspace" aria-label="Pathfinder workspace">
			<Controls
				{algorithm}
				{speed}
				{isRunning}
				{isPaused}
				{status}
				onAlgorithmChange={(value) => (algorithm = value)}
				onSpeedChange={(value) => (speed = value)}
				onStart={startSearch}
				onPauseToggle={togglePause}
				onClearPath={clearPath}
				onClearWalls={clearWalls}
				onReset={resetGrid}
        onStop={stopSearch}
			/>

      <Stats {algorithm} {stats} />

			<div class="grid-frame">
				<Grid
					rows={GRID_ROWS}
					cols={GRID_COLS}
					{start}
					{goal}
					{walls}
					{explored}
					{path}
					disabled={isRunning}
					onWallChange={updateWall}
          onEndpointChange={updateEndpoint}
				/>
			</div>
		</section>
	</main>
</div>

<style>
	:global(*) {
		box-sizing: border-box;
	}

	:global(html) {
		min-width: 320px;
		background: #07101e;
	}

	:global(body) {
		min-width: 320px;
		min-height: 100vh;
		margin: 0;
		background:
			radial-gradient(circle at 15% 0%, #15315a 0, transparent 34rem),
			#07101e;
		color: #edf3ff;
		font-family:
			Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
			"Segoe UI", sans-serif;
	}

	:global(button),
	:global(select),
	:global(input) {
		font: inherit;
	}

	.page-shell {
		width: min(92rem, 100%);
		margin: 0 auto;
		padding: clamp(1.25rem, 3vw, 3rem);
	}

	.hero {
		display: flex;
		align-items: end;
		justify-content: space-between;
		gap: 2rem;
		margin-bottom: 1.5rem;
	}

	.eyebrow {
		margin: 0 0 0.4rem;
		color: #66d8ff;
		font-size: 0.75rem;
		font-weight: 750;
		letter-spacing: 0.14em;
		text-transform: uppercase;
	}

	h1 {
		margin: 0;
		font-size: clamp(2.5rem, 7vw, 5.5rem);
		line-height: 0.95;
		letter-spacing: -0.065em;
	}

	.introduction {
		max-width: 38rem;
		margin: 1rem 0 0;
		color: #aebbd0;
		font-size: clamp(1rem, 2vw, 1.15rem);
		line-height: 1.6;
	}

	.legend {
		display: flex;
		flex-wrap: wrap;
		justify-content: flex-end;
		gap: 0.65rem 1rem;
		color: #aebbd0;
		font-size: 0.8rem;
	}

	.legend span {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		white-space: nowrap;
	}

	.legend i {
		width: 0.72rem;
		height: 0.72rem;
		border-radius: 0.2rem;
	}

	.start-swatch {
		background: #24a76a;
	}

	.goal-swatch {
		background: #e55353;
	}

	.wall-swatch {
		background: #546178;
	}

	.visited-swatch {
		background: #185b78;
	}

	.path-swatch {
		background: #f3c745;
	}

	.workspace {
		padding: clamp(0.75rem, 2vw, 1.25rem);
		border: 1px solid #25344d;
		border-radius: 1.25rem;
		background: rgb(11 22 39 / 88%);
		box-shadow: 0 2rem 5rem rgb(0 0 0 / 24%);
		backdrop-filter: blur(1rem);
	}

	.grid-frame {
		overflow-x: auto;
	}

	.grid-frame :global(.grid) {
		min-width: 44rem;
	}

	@media (max-width: 50rem) {
		.hero {
			display: grid;
		}

		.legend {
			justify-content: flex-start;
		}
	}
</style>

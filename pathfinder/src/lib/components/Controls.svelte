<script lang="ts">
	import type { Algorithm } from '$lib/types';

	type ControlsProps = {
		algorithm: Algorithm;
		speed: number;
		isRunning: boolean;
		isPaused: boolean;
		status: string;
		onAlgorithmChange: (algorithm: Algorithm) => void;
		onSpeedChange: (speed: number) => void;
		onStart: () => void;
		onPauseToggle: () => void;
		onClearPath: () => void;
		onClearWalls: () => void;
		onReset: () => void;
    onStop: () => void;
	};

	let {
		algorithm,
		speed,
		isRunning,
		isPaused,
		status,
		onAlgorithmChange,
		onSpeedChange,
		onStart,
		onPauseToggle,
		onClearPath,
		onClearWalls,
		onReset,
    onStop
	}: ControlsProps = $props();

	function handleAlgorithmChange(event: Event): void {
		const select = event.currentTarget as HTMLSelectElement;
		onAlgorithmChange(select.value as Algorithm);
	}

	function handleSpeedChange(event: Event): void {
		const input = event.currentTarget as HTMLInputElement;
		onSpeedChange(Number(input.value));
	}
</script>

<div class="controls">
	<div class="field">
		<label for="algorithm">Algorithm</label>
		<select
			id="algorithm"
			value={algorithm}
			disabled={isRunning}
			onchange={handleAlgorithmChange}
		>
			<option value="astar">A* search</option>
			<option value="dijkstra">Dijkstra</option>
			<option value="bfs">Breadth-first search</option>
		</select>
	</div>

	<div class="field speed-field">
		<div class="label-row">
			<label for="speed">Animation</label>
			<output for="speed">{speed} ms</output>
		</div>

		<input
			id="speed"
			type="range"
			min="5"
			max="100"
			step="5"
			value={speed}
			oninput={handleSpeedChange}
		/>
	</div>

	<div class="actions">
		<button
			class="primary"
			type="button"
			disabled={isRunning}
			onclick={onStart}
		>
			Run search
		</button>

		<button
			type="button"
			disabled={!isRunning}
			onclick={onPauseToggle}
		>
			{isPaused ? 'Resume' : 'Pause'}
		</button>

    <button
      type="button"
      disabled={!isRunning}
      onclick={onStop}
    >
      Stop
    </button>

		<button
			type="button"
			disabled={isRunning}
			onclick={onClearPath}
		>
			Clear path
		</button>

		<button
			type="button"
			disabled={isRunning}
			onclick={onClearWalls}
		>
			Clear walls
		</button>

		<button
			type="button"
			disabled={isRunning}
			onclick={onReset}
		>
			Reset
		</button>
	</div>

	<p class="status" aria-live="polite">
		<span class:active={isRunning}></span>
		{status}
	</p>
</div>

<style>
	.controls {
		display: grid;
		grid-template-columns: minmax(10rem, 0.65fr) minmax(12rem, 1fr) auto;
		align-items: end;
		gap: 1rem;
		margin-bottom: 1rem;
		padding: 1rem;
		border: 1px solid #293a56;
		border-radius: 0.85rem;
		background: #0d192c;
	}

	.field {
		display: grid;
		gap: 0.45rem;
	}

	label,
	.label-row {
		color: #aebbd0;
		font-size: 0.76rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
	}

	.label-row {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
	}

	output {
		color: #66d8ff;
		font-variant-numeric: tabular-nums;
	}

	select {
		min-height: 2.6rem;
		padding: 0 2.25rem 0 0.75rem;
		border: 1px solid #3b4d69;
		border-radius: 0.55rem;
		background: #142137;
		color: #edf3ff;
	}

	input[type="range"] {
		width: 100%;
		accent-color: #66d8ff;
	}

	.actions {
		display: flex;
		flex-wrap: wrap;
		justify-content: flex-end;
		gap: 0.5rem;
	}

	button {
		min-height: 2.6rem;
		padding: 0 0.85rem;
		border: 1px solid #3b4d69;
		border-radius: 0.55rem;
		background: #142137;
		color: #dce7f8;
		cursor: pointer;
	}

	button:hover:not(:disabled) {
		border-color: #657b9d;
		background: #1b2a43;
	}

	button.primary {
		border-color: #45badf;
		background: #1681a4;
		color: white;
		font-weight: 750;
	}

	button:focus-visible,
	select:focus-visible,
	input:focus-visible {
		outline: 2px solid #79ddff;
		outline-offset: 2px;
	}

	button:disabled,
	select:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}

	.status {
		display: flex;
		grid-column: 1 / -1;
		align-items: center;
		gap: 0.5rem;
		margin: 0;
		color: #8f9db3;
		font-size: 0.85rem;
	}

	.status span {
		width: 0.55rem;
		height: 0.55rem;
		border-radius: 50%;
		background: #526078;
	}

	.status span.active {
		background: #66d8ff;
		box-shadow: 0 0 0.8rem #66d8ff;
	}

	@media (max-width: 64rem) {
		.controls {
			grid-template-columns: 1fr 1fr;
		}

		.actions {
			grid-column: 1 / -1;
			justify-content: flex-start;
		}
	}

	@media (max-width: 42rem) {
		.controls {
			grid-template-columns: 1fr;
		}

		.actions {
			grid-column: auto;
		}

		.status {
			grid-column: auto;
		}
	}
</style>

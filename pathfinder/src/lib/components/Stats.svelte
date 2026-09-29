<script lang="ts">
	import type { Algorithm, SearchStats } from '$lib/types';

	type StatsProps = {
		algorithm: Algorithm;
		stats: SearchStats | null;
	};

	let { algorithm, stats }: StatsProps = $props();

	function algorithmLabel(value: Algorithm): string {
		switch (value) {
			case 'astar':
				return 'A*';
			case 'dijkstra':
				return 'Dijkstra';
			case 'bfs':
				return 'BFS';
		}
	}

	function durationLabel(milliseconds: number): string {
		if (milliseconds < 1) {
			return `${milliseconds.toFixed(3)} ms`;
		}

		return `${milliseconds.toFixed(2)} ms`;
	}
</script>

<div class="stats" aria-label="Search statistics">
	<div class="stat">
		<span>Algorithm</span>
		<strong>{algorithmLabel(algorithm)}</strong>
	</div>

	<div class="stat">
		<span>Explored</span>
		<strong>{stats?.visited ?? '—'}</strong>
	</div>

	<div class="stat">
		<span>Path length</span>
		<strong>{stats?.found ? stats.pathLength : '—'}</strong>
	</div>

	<div class="stat">
		<span>Backend time</span>
		<strong>
			{stats ? durationLabel(stats.durationMs) : '—'}
		</strong>
	</div>

	<div class="stat">
		<span>Result</span>
		<strong class:success={stats?.found} class:failure={stats && !stats.found}>
			{stats ? (stats.found ? 'Found' : 'No path') : 'Ready'}
		</strong>
	</div>
</div>

<style>
	.stats {
		display: grid;
		grid-template-columns: repeat(5, minmax(0, 1fr));
		margin-bottom: 1rem;
		overflow: hidden;
		border: 1px solid #293a56;
		border-radius: 0.85rem;
		background: #0d192c;
	}

	.stat {
		display: grid;
		gap: 0.35rem;
		min-width: 0;
		padding: 0.85rem 1rem;
		border-right: 1px solid #293a56;
	}

	.stat:last-child {
		border-right: 0;
	}

	.stat span {
		overflow: hidden;
		color: #77869d;
		font-size: 0.68rem;
		font-weight: 750;
		letter-spacing: 0.08em;
		text-overflow: ellipsis;
		text-transform: uppercase;
		white-space: nowrap;
	}

	.stat strong {
		overflow: hidden;
		color: #edf3ff;
		font-size: 1rem;
		font-variant-numeric: tabular-nums;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.stat strong.success {
		color: #5be09a;
	}

	.stat strong.failure {
		color: #ff8585;
	}

	@media (max-width: 48rem) {
		.stats {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}

		.stat {
			border-bottom: 1px solid #293a56;
		}

		.stat:nth-child(2n) {
			border-right: 0;
		}

		.stat:last-child {
			border-bottom: 0;
		}
	}
</style>

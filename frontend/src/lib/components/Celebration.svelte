<script lang="ts">
	import { fly, fade } from 'svelte/transition';

	let {
		winner = '',
		winnerNumero = 0,
		autoNombres = {} as Record<string, number>,
		show = false,
		onclose
	} = $props();

	let confetti = $state<Array<{id: number; x: number; delay: number; color: string; size: number}>>([]);
	let confettiInterval: ReturnType<typeof setInterval> | undefined;

	$effect(() => {
		if (show) {
			// Generar confetti
			confetti = Array.from({ length: 50 }, (_, i) => ({
				id: i,
				x: Math.random() * 100,
				delay: Math.random() * 2,
				color: ['#f59e0b', '#ef4444', '#059669', '#3b82f6', '#8b5cf6', '#ec4899'][Math.floor(Math.random() * 6)],
				size: Math.random() * 8 + 4
			}));
		} else {
			confetti = [];
		}
	});

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose?.();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
	<div class="celebration-overlay" onclick={onclose} role="presentation">
		<!-- Confetti -->
		{#each confetti as c (c.id)}
			<div
				class="confetti-piece"
				style="left: {c.x}%; animation-delay: {c.delay}s; background: {c.color}; width: {c.size}px; height: {c.size * 0.6}px;"
				in:fade={{ duration: 100, delay: c.delay * 500 }}
			></div>
		{/each}

		<div class="celebration-content" onclick={(e: MouseEvent) => e.stopPropagation()}>
			<div class="trophy" in:fly={{ y: -50, duration: 600 }}>🏆</div>
			<h1 class="celebration-title" in:fly={{ y: 30, duration: 500, delay: 300 }}>
				¡Campeón!
			</h1>
			<div class="winner-display" in:fly={{ y: 30, duration: 500, delay: 500 }}>
				<span class="winner-numero">#{winnerNumero}</span>
				<span class="winner-name">{winner}</span>
			</div>
			<p class="celebration-sub" in:fade={{ duration: 500, delay: 800 }}>
				Ganador de la Gran Final — Pinewood Derby D15
			</p>
			<button class="btn-close" onclick={onclose} in:fade={{ delay: 1000 }}>
				Cerrar
			</button>
		</div>
	</div>
{/if}

<style>
	.celebration-overlay {
		position: fixed;
		inset: 0;
		background: radial-gradient(ellipse at center, #1a1f2e 0%, #0a0f1a 70%);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 300;
		overflow: hidden;
	}

	@keyframes confetti-fall {
		0% { transform: translateY(-10vh) rotate(0deg); opacity: 1; }
		100% { transform: translateY(110vh) rotate(720deg); opacity: 0; }
	}

	.confetti-piece {
		position: absolute;
		top: -10vh;
		border-radius: 2px;
		animation: confetti-fall 4s ease-in forwards;
	}

	.celebration-content {
		text-align: center;
		z-index: 1;
		padding: 2rem;
	}

	.trophy {
		font-size: 5rem;
		margin-bottom: 1rem;
		filter: drop-shadow(0 0 30px rgba(255,215,0,0.5));
	}

	.celebration-title {
		color: var(--racing-amber);
		font-size: 4rem;
		margin: 0 0 1rem 0;
		text-shadow: 0 0 30px rgba(245,158,11,0.5);
	}

	.winner-display {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.winner-numero {
		font-size: 3rem;
		font-weight: 900;
		color: var(--racing-amber);
	}

	.winner-name {
		font-size: 2.5rem;
		font-weight: 700;
		color: white;
	}

	.celebration-sub {
		color: var(--racing-text-dim);
		font-size: 1.1rem;
		margin-bottom: 2rem;
	}

	.btn-close {
		padding: 0.75rem 2rem;
		background: var(--racing-amber);
		color: var(--racing-black);
		border: none;
		border-radius: 0.25rem;
		font-weight: 700;
		font-size: 1rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-close:hover {
		transform: scale(1.05);
		box-shadow: 0 0 20px rgba(245,158,11,0.4);
	}

	@media (max-width: 640px) {
		.celebration-title { font-size: 2.5rem; }
		.winner-name { font-size: 1.5rem; }
		.winner-numero { font-size: 2rem; }
		.trophy { font-size: 3.5rem; }
	}
</style>

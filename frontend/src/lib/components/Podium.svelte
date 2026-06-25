<script lang="ts">
	import { fly, scale } from 'svelte/transition';

	let {
		ordenLlegada = [] as string[],
		autoNombres = {} as Record<string, string>,
		autoNumeros = {} as Record<string, number>,
		label = 'Podio',
		show = false,
		onclose
	} = $props();

	let revealed = $state(0);

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose?.();
	}

	$effect(() => {
		if (show) {
			revealed = 0;
			const timer = setInterval(() => {
				revealed++;
				if (revealed >= ordenLlegada.length) clearInterval(timer);
			}, 400);
			return () => clearInterval(timer);
		} else {
			revealed = 0;
		}
	});

	const medallas = ['🥇', '🥈', '🥉', '4️⃣'];
	const colores = ['gold', 'silver', '#cd7f32', 'transparent'];
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show && ordenLlegada.length > 0}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="podium-overlay" onclick={onclose} role="presentation">
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<div class="podium-content" onclick={(e: MouseEvent) => e.stopPropagation()}>
			<h2 class="podium-title" in:fly={{ y: -20, duration: 300 }}>{label}</h2>

			<div class="podium-steps">
				{#each ordenLlegada as autoId, i}
					{#if revealed > i}
						<div
							class="podium-step"
							style="--rank: {i}; --color: {colores[i]}"
							in:scale={{ duration: 300, delay: i * 100 }}
						>
							<div class="medal">{medallas[i] || '#' + (i + 1)}</div>
							<div class="podium-bar" style="height: {100 - i * 20}px; background: {colores[i]}">
								<span class="podium-pos">{i + 1}°</span>
							</div>
							<div class="podium-name">
								<span class="numero">#{autoNumeros[autoId]}</span>
								{autoNombres[autoId] || autoId}
							</div>
						</div>
					{/if}
				{/each}
			</div>
		</div>
	</div>
{/if}

<style>
	.podium-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 200;
		backdrop-filter: blur(4px);
	}

	.podium-content {
		text-align: center;
		padding: 2rem;
	}

	.podium-title {
		color: var(--orange);
		font-size: 2.5rem;
		margin: 0 0 2rem 0;
		text-shadow: 0 0 30px rgba(245,158,11,0.5);
	}

	.podium-steps {
		display: flex;
		justify-content: center;
		align-items: flex-end;
		gap: 1.5rem;
	}

	.podium-step {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.medal {
		font-size: 2.5rem;
		filter: drop-shadow(0 0 10px rgba(255,215,0,0.5));
	}

	.podium-bar {
		width: 80px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 0.25rem 0.25rem 0 0;
		min-height: 40px;
		transition: all 0.3s;
		box-shadow: 0 -4px 15px rgba(0,0,0,0.3);
	}

	.podium-pos {
		color: var(--arcade-black);
		font-weight: 900;
		font-size: 1.5rem;
	}

	.podium-name {
		color: var(--text-primary);
		font-weight: 600;
		font-size: 0.9rem;
		text-align: center;
		max-width: 120px;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.numero {
		color: var(--orange);
		margin-right: 0.25rem;
	}

	@media (max-width: 640px) {
		.podium-title { font-size: 1.5rem; }
		.podium-bar { width: 60px; }
		.medal { font-size: 1.8rem; }
	}
</style>

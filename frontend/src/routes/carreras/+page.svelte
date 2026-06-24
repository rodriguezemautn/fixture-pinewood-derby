<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	let categorias = $state<any[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const res = await fetch('/api/categorias');
			if (res.ok) categorias = await res.json();
		} catch {} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Carreras — Fixture D15</title>
</svelte:head>

<div class="page">
	<header class="page-header">
		<div class="header-inner">
			<img src="/assets/derby_logo.jpg" alt="Derby" class="logo" />
			<div>
				<h1>Carreras</h1>
				<p class="subtitle">Pinewood Derby — Destacamento 15</p>
			</div>
		</div>
		<div class="header-stripe"></div>
	</header>

	<main class="main-content">
		{#if loading}
			<p class="loading">Cargando categorías...</p>
		{:else if categorias.length === 0}
			<div class="empty" in:fade>
				<img src="/assets/emblema.jpg" alt="Emblema" class="empty-img" />
				<h2>No hay carreras disponibles</h2>
				<p>El administrador debe crear categorías y generar el fixture</p>
			</div>
		{:else}
			<div class="categorias-grid">
				{#each categorias as cat (cat.id)}
					<a href="/carreras/{cat.id}" class="categoria-card" in:fly={{ y: 20, duration: 300 }}>
						<span class="cat-nombre">{cat.nombre}</span>
						<span class="cat-edades">{cat.edad_min} - {cat.edad_max} años</span>
						<span class="cat-cta">👁️ Ver Carreras</span>
					</a>
				{/each}
			</div>
		{/if}
	</main>
</div>

<style>
	.page { min-height: 100vh; background: var(--racing-black); }

	.page-header {
		background: linear-gradient(180deg, #0a0f1a, #0f172a);
		padding: 1rem;
	}

	.header-inner {
		display: flex; align-items: center; gap: 1rem;
		max-width: 800px; margin: 0 auto;
	}

	.logo { height: 48px; width: auto; border-radius: 0.25rem; }
	h1 { color: var(--racing-amber); font-size: 1.5rem; margin: 0; }
	.subtitle { color: var(--racing-text-dim); font-size: 0.8rem; margin: 0; }

	.header-stripe {
		height: 3px; margin-top: 0.75rem;
		background: repeating-linear-gradient(90deg, var(--racing-amber) 0, var(--racing-amber) 8px, transparent 8px, transparent 16px);
	}

	.main-content { max-width: 800px; margin: 2rem auto; padding: 0 1rem; }

	.loading { text-align: center; color: var(--racing-text-dim); padding: 3rem; }

	.empty { text-align: center; padding: 4rem 1rem; }
	.empty-img { width: 100px; border-radius: 50%; opacity: 0.4; margin-bottom: 1rem; }
	.empty h2 { color: var(--racing-text-dim); }
	.empty p { color: var(--racing-text-dim); font-size: 0.9rem; }

	.categorias-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 1rem;
	}

	.categoria-card {
		display: flex; flex-direction: column; gap: 0.25rem;
		background: var(--racing-gray); border: 1px solid var(--racing-border);
		border-radius: 0.5rem; padding: 1.5rem;
		text-decoration: none; transition: all 0.2s;
	}

	.categoria-card:hover {
		border-color: var(--racing-amber);
		transform: translateY(-3px);
		box-shadow: 0 8px 24px rgba(0,0,0,0.3);
	}

	.cat-nombre { color: var(--racing-amber); font-weight: 700; font-size: 1.1rem; }
	.cat-edades { color: var(--racing-text-dim); font-size: 0.85rem; }
	.cat-cta { margin-top: 0.5rem; color: var(--racing-amber); font-size: 0.85rem; font-weight: 600; }

	@media (max-width: 640px) { .categorias-grid { grid-template-columns: 1fr; } }
</style>

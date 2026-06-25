<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	let autos = $state<any[]>([]);
	let categorias = $state<any[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const [autosRes, catsRes] = await Promise.all([
				fetch('/api/autos'),
				fetch('/api/categorias')
			]);
			if (autosRes.ok) autos = await autosRes.json();
			if (catsRes.ok) categorias = await catsRes.json();
		} catch {} finally {
			loading = false;
		}
	});

	// Armar mapa de categoria_id → nombre
	const catMap = $derived(() => {
		const m: Record<string, string> = {};
		for (const c of categorias) m[c.id] = c.nombre;
		return m;
	});

	// Agrupar autos por categoría
	const autosPorCategoria = $derived.by(() => {
		const grupos: Record<string, any[]> = {};
		for (const a of autos) {
			const catId = a.categoria_id;
			if (!grupos[catId]) grupos[catId] = [];
			grupos[catId].push(a);
		}
		return grupos;
	});
</script>

<svelte:head>
	<title>Autos — Fixture D15</title>
</svelte:head>

<div class="page">
	<header class="page-header">
		<div class="header-inner">
			<a href="/carreras" class="back-link">← Carreras</a>
			<img src="/assets/derby_logo.jpg" alt="Derby" class="logo" />
			<div>
				<h1>Autos Registrados</h1>
				<p class="subtitle">Pinewood Derby — Destacamento 15</p>
			</div>
		</div>
		<div class="header-stripe"></div>
	</header>

	<main class="main-content">
		{#if loading}
			<p class="loading">Cargando autos...</p>
		{:else if autos.length === 0}
			<div class="empty" in:fade>
				<img src="/assets/auto_madera.jpg" alt="Autos" class="empty-img" />
				<h2>No hay autos registrados</h2>
				<p>El administrador debe cargar los autos de cada categoría</p>
			</div>
		{:else}
			{#each Object.entries(autosPorCategoria) as [catId, autosCat]}
				<section class="cat-section" in:fly={{ y: 15, duration: 300 }}>
					<h2 class="cat-title">{catMap()[catId] || catId}</h2>
					<div class="autos-grid">
						{#each autosCat as auto (auto.id)}
							<div class="auto-card">
								{#if auto.foto_url}
									<img src={auto.foto_url} alt={auto.nombre} class="auto-foto" />
								{:else}
									<div class="auto-foto-placeholder">🏎️</div>
								{/if}
								<div class="auto-info">
									<span class="auto-numero">#{auto.numero}</span>
									<span class="auto-nombre">{auto.nombre}</span>
								<span class="auto-creador">por {auto.creador}</span>
								<span class="auto-edad">{auto.edad} años</span>
								{#if auto.peso > 0}<span class="auto-peso">{auto.peso}g</span>{/if}
								</div>
							</div>
						{/each}
					</div>
				</section>
			{/each}

			<p class="total">{autos.length} auto{autos.length !== 1 ? 's' : ''} registrado{autos.length !== 1 ? 's' : ''}</p>
		{/if}
	</main>
</div>

<style>
	.page { min-height: 100vh; background: var(--arcade-black); }

	.page-header { background: linear-gradient(180deg, var(--arcade-black), var(--arcade-dark)); padding: 1rem; }

	.header-inner {
		display: flex; align-items: center; gap: 1rem;
		max-width: 800px; margin: 0 auto; flex-wrap: wrap;
	}

	.back-link { color: var(--text-dim); text-decoration: none; font-size: 0.85rem; width: 100%; }
	.back-link:hover { color: var(--orange); }

	.logo { height: 48px; width: auto; border-radius: 0.25rem; }
	h1 { color: var(--orange); font-size: 1.5rem; margin: 0; }
	.subtitle { color: var(--text-dim); font-size: 0.8rem; margin: 0; }
	.header-stripe { height: 3px; margin-top: 0.75rem; background: repeating-linear-gradient(90deg, var(--orange) 0, var(--orange) 8px, transparent 8px, transparent 16px); }

	.main-content { max-width: 800px; margin: 2rem auto; padding: 0 1rem; }

	.loading { text-align: center; color: var(--text-dim); padding: 3rem; }

	.empty { text-align: center; padding: 4rem 1rem; }
	.empty-img { width: 100px; border-radius: 0.5rem; opacity: 0.4; margin-bottom: 1rem; }
	.empty h2 { color: var(--text-dim); font-size: 1.2rem; }
	.empty p { color: var(--text-dim); font-size: 0.9rem; }

	.cat-section { margin-bottom: 2rem; }
	.cat-title { color: var(--orange); font-size: 1.1rem; margin: 0 0 0.75rem 0; padding-bottom: 0.5rem; border-bottom: 1px solid var(--border-color); }

	.autos-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 0.75rem;
	}

	.auto-card {
		display: flex; align-items: center; gap: 0.75rem;
		background: var(--arcade-dark); border: 1px solid var(--border-color);
		border-radius: 0.5rem; padding: 0.75rem;
		transition: border-color 0.2s;
	}

	.auto-card:hover { border-color: var(--orange); }

	.auto-foto { width: 52px; height: 52px; object-fit: cover; border-radius: 0.25rem; border: 1px solid var(--border-color); }
	.auto-foto-placeholder { width: 52px; height: 52px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; background: var(--arcade-surface); border-radius: 0.25rem; }

	.auto-info { display: flex; flex-direction: column; gap: 0.1rem; min-width: 0; }
	.auto-numero { color: var(--orange); font-weight: 700; font-size: 1rem; }
	.auto-nombre { color: var(--text-primary); font-size: 0.9rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
	.auto-creador { color: var(--text-dim); font-size: 0.8rem; }
	.auto-edad { color: var(--text-dim); font-size: 0.75rem; }
	.auto-peso { color: var(--orange); font-size: 0.75rem; font-weight: 600; }

	.total { text-align: center; color: var(--text-dim); font-size: 0.85rem; padding: 1rem; }

	@media (max-width: 640px) { .autos-grid { grid-template-columns: 1fr; } }
</style>

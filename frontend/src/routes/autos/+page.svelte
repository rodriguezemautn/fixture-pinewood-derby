<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly, scale } from 'svelte/transition';

	let autos = $state<any[]>([]);
	let categorias = $state<any[]>([]);
	let loading = $state(true);

	let selectedAuto = $state<any | null>(null);

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

	function openDetail(auto: any) {
		selectedAuto = auto;
	}

	function closeDetail() {
		selectedAuto = null;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') closeDetail();
	}

	const categoriaNombre = $derived(
		selectedAuto ? catMap()[selectedAuto.categoria_id] || '—' : ''
	);
</script>

<svelte:window onkeydown={handleKeydown} />

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
							<button class="auto-card" onclick={() => openDetail(auto)}>
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
							</button>
						{/each}
					</div>
				</section>
			{/each}

			<p class="total">{autos.length} auto{autos.length !== 1 ? 's' : ''} registrado{autos.length !== 1 ? 's' : ''}</p>
		{/if}
	</main>

	<!-- Modal de detalle -->
	{#if selectedAuto}
		<!-- svelte-ignore a11y_click_events_have_key_events a11y_interactive_supports_focus -->
		<div class="modal-overlay" onclick={closeDetail} role="presentation" in:fade={{ duration: 150 }}>
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<div class="modal" onclick={(e) => e.stopPropagation()} role="dialog" aria-label="Detalle del auto" in:scale={{ duration: 200, start: 0.9 }}>
				<button class="modal-close" onclick={closeDetail}>✕</button>

				<div class="modal-foto">
					{#if selectedAuto.foto_url}
						<img src={selectedAuto.foto_url} alt={selectedAuto.nombre} class="modal-img" />
					{:else}
						<div class="modal-img-placeholder">🏎️</div>
					{/if}
				</div>

				<div class="modal-body">
					<h2 class="modal-numero">#{selectedAuto.numero}</h2>
					<h3 class="modal-nombre">{selectedAuto.nombre}</h3>

					<div class="modal-details">
						<div class="detail-row">
							<span class="detail-label">Creador</span>
							<span class="detail-value">{selectedAuto.creador}</span>
						</div>
						<div class="detail-row">
							<span class="detail-label">Edad</span>
							<span class="detail-value">{selectedAuto.edad} años</span>
						</div>
						{#if selectedAuto.peso > 0}
							<div class="detail-row">
								<span class="detail-label">Peso</span>
								<span class="detail-value">{selectedAuto.peso}g</span>
							</div>
						{/if}
						<div class="detail-row">
							<span class="detail-label">Categoría</span>
							<span class="detail-value">{categoriaNombre}</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
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
		transition: border-color 0.2s, transform 0.15s;
		cursor: pointer; text-align: left; width: 100%;
		font-family: inherit; font-size: inherit; color: inherit;
		appearance: none; -webkit-appearance: none;
	}

	.auto-card:hover { border-color: var(--orange); transform: translateY(-2px); }
	.auto-card:focus-visible { outline: 2px solid var(--orange); outline-offset: 2px; }

	.auto-foto { width: 52px; height: 52px; object-fit: cover; border-radius: 0.25rem; border: 1px solid var(--border-color); }
	.auto-foto-placeholder { width: 52px; height: 52px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; background: var(--arcade-surface); border-radius: 0.25rem; }

	.auto-info { display: flex; flex-direction: column; gap: 0.1rem; min-width: 0; }
	.auto-numero { color: var(--orange); font-weight: 700; font-size: 1rem; }
	.auto-nombre { color: var(--text-primary); font-size: 0.9rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
	.auto-creador { color: var(--text-dim); font-size: 0.8rem; }
	.auto-edad { color: var(--text-dim); font-size: 0.75rem; }
	.auto-peso { color: var(--orange); font-size: 0.75rem; font-weight: 600; }

	.total { text-align: center; color: var(--text-dim); font-size: 0.85rem; padding: 1rem; }

	/* ─── Modal ─── */
	.modal-overlay {
		position: fixed; inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex; align-items: center; justify-content: center;
		z-index: 200;
		padding: 1rem;
	}

	.modal {
		background: var(--arcade-surface);
		border: 2px solid var(--orange-border);
		border-radius: 0.75rem;
		width: 100%;
		max-width: 420px;
		max-height: 90vh;
		overflow-y: auto;
		box-shadow: 0 8px 32px rgba(0,0,0,0.5);
		position: relative;
	}

	.modal-close {
		position: absolute; top: 0.5rem; right: 0.5rem;
		background: var(--arcade-dark); border: 1px solid var(--border-color);
		color: var(--text-dim); width: 32px; height: 32px;
		border-radius: 50%;
		font-size: 1rem; cursor: pointer;
		display: flex; align-items: center; justify-content: center;
		z-index: 10; transition: all 0.15s;
	}
	.modal-close:hover { background: var(--orange); color: var(--arcade-black); }

	.modal-foto { width: 100%; height: 260px; overflow: hidden; border-radius: 0.75rem 0.75rem 0 0; background: var(--arcade-dark); }
	.modal-img { width: 100%; height: 100%; object-fit: cover; }
	.modal-img-placeholder {
		width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;
		font-size: 4rem; color: var(--text-dim);
	}

	.modal-body { padding: 1.25rem; }
	.modal-numero { color: var(--orange); font-size: 1.5rem; margin: 0; }
	.modal-nombre { color: var(--text-primary); font-size: 1.1rem; margin: 0.25rem 0 1rem 0; }

	.modal-details { display: flex; flex-direction: column; gap: 0.5rem; }
	.detail-row {
		display: flex; justify-content: space-between; align-items: center;
		padding: 0.5rem 0; border-bottom: 1px solid var(--border-color);
	}
	.detail-row:last-child { border-bottom: none; }
	.detail-label { color: var(--text-dim); font-size: 0.85rem; }
	.detail-value { color: var(--text-primary); font-weight: 600; font-size: 0.95rem; }

	/* ─── Responsive ─── */
	@media (max-width: 640px) {
		.autos-grid { grid-template-columns: 1fr; }
		.modal-foto { height: 200px; }
	}
</style>

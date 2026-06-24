<script lang="ts">
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';
	import { browser } from '$app/environment';

	let loaded = $state(false);
	let categorias = $state<any[]>([]);
	let fixturesStatus = $state<Record<string, string>>({});
	let loadingCats = $state(true);

	$effect(() => {
		loaded = true;
	});

	onMount(() => {
		cargarCategorias();
	});

	async function cargarCategorias() {
		try {
			const res = await fetch('/api/categorias');
			if (res.ok) {
				const cats = await res.json();
				categorias = cats;
				// Verificar estado de fixture de cada categoría
				for (const cat of cats) {
					fetch(`/api/categorias/${cat.id}/fixture`).then(r => {
						if (r.ok) r.json().then(f => {
							if (f && f.estado) fixturesStatus[cat.id] = f.estado;
						});
					}).catch(() => {});
				}
			}
		} catch {} finally {
			loadingCats = false;
		}
	}

	const isAutenticado = $derived(browser ? !!localStorage.getItem('auth_token') : false);

	function estadoCat(catId: string): string {
		return fixturesStatus[catId] || '';
	}
</script>

<svelte:head>
	<title>Fixture — Pinewood Derby D15</title>
</svelte:head>

<!-- ════ HERO ════ -->
<section class="hero">
	<div class="hero-bg"></div>
	<div class="hero-speed"></div>
	<div class="hero-content">
		{#if loaded}
			<div class="hero-logo" in:fly={{ y: -30, duration: 600 }}>
				<img src="/assets/emblema.jpg" alt="Emblema D15" class="emblema" />
			</div>
			<div class="hero-text" in:fly={{ y: 30, duration: 600, delay: 200 }}>
				<h1 class="hero-title">
					<span class="title-line">Pinewood Derby</span>
					<span class="title-line subtitle">Destacamento 15</span>
				</h1>
				<p class="hero-desc">Controlador de fixture para carreras de autitos de madera</p>
				<div class="hero-actions">
					{#if isAutenticado}
						<a href="/admin/categorias" class="btn-racing-primary">🏎️ Gestionar Carreras</a>
					{:else}
						<a href="/login" class="btn-racing-primary">🔐 Admin</a>
					{/if}
				</div>
			</div>
		{/if}
	</div>
	<div class="hero-stripe"></div>
</section>

<!-- ════ CARRERAS DISPONIBLES ════ -->
<section class="info-section">
	{#if loaded}
		<div class="section-title" in:fade>
			<h2>🏁 Carreras</h2>
			<p class="section-desc">Seleccioná una categoría para ver las carreras en vivo</p>
		</div>

		{#if loadingCats}
			<div class="categorias-grid">
				<div class="skeleton-card"></div>
				<div class="skeleton-card"></div>
				<div class="skeleton-card"></div>
			</div>
		{:else if categorias.length === 0}
			<div class="empty-carreras" in:fade>
				<img src="/assets/emblema.jpg" alt="Emblema" class="empty-img" />
				<p>No hay carreras disponibles todavía</p>
				<p class="empty-sub">El administrador debe crear categorías y generar el fixture</p>
			</div>
		{:else}
			<div class="categorias-grid">
				{#each categorias as cat (cat.id)}
					<a href="/carreras/{cat.id}" class="categoria-card" class:live={estadoCat(cat.id) === 'en_curso'} class:finalizado={estadoCat(cat.id) === 'finalizado'} in:fly={{ y: 20, duration: 400 }}>
						<span class="cat-nombre">{cat.nombre}</span>
						<span class="cat-edades">{cat.edad_min} - {cat.edad_max} años</span>
						<span class="cat-badge-row">
							{#if estadoCat(cat.id) === 'en_curso'}
								<span class="badge-live">🔴 En Vivo</span>
							{:else if estadoCat(cat.id) === 'finalizado'}
								<span class="badge-done">🏁 Finalizado</span>
							{/if}
							<span class="cat-cta">👁️ Ver Carreras →</span>
						</span>
					</a>
				{/each}
			</div>
		{/if}
	{/if}
</section>

<!-- ════ INFO CARDS ════ -->
<section class="info-section cards-section">
	{#if loaded}
		<div class="section-title" in:fade>
			<h2>¿Cómo funciona?</h2>
		</div>
		<div class="info-grid" in:fade={{ duration: 500, delay: 400 }}>
			<a href="/carreras" class="info-card" in:fly={{ y: 20, duration: 400, delay: 500 }}>
				<img src="/assets/auto_madera.jpg" alt="Auto de madera" class="card-img" />
				<h3>Categorías</h3>
				<p>Explorá las categorías disponibles y seguí las carreras en vivo</p>
				<span class="card-link">👁️ Ver carreras →</span>
			</a>
			<a href="/autos" class="info-card" in:fly={{ y: 20, duration: 400, delay: 600 }}>
				<img src="/assets/D15_logo.jpg" alt="D15 Logo" class="card-img" />
				<h3>Autos</h3>
				<p>Explorá todos los autos registrados por categoría</p>
				<span class="card-link">🏎️ Ver autos →</span>
			</a>
			<a href="/carreras" class="info-card" in:fly={{ y: 20, duration: 400, delay: 700 }}>
				<img src="/assets/derbi.jpg" alt="Derbi" class="card-img" />
				<h3>Carreras</h3>
				<p>Seguí las carreras en vivo con posiciones, podios y resultados</p>
				<span class="card-link">🏁 Ver carreras →</span>
			</a>
		</div>
	{/if}
</section>

<!-- ════ SPONSORS / CREDITS ════ -->
<section class="sponsors">
	<div class="sponsor-logos">
		<img src="/assets/logo.jpg" alt="Logo" class="sponsor-img" />
		<img src="/assets/logo_naranja.jpg" alt="Logo Naranja" class="sponsor-img" />
		<img src="/assets/D15_logo.jpg" alt="D15" class="sponsor-img" />
	</div>
</section>

<style>
	/* ─── Hero ─────────────────────────────────── */
	.hero {
		position: relative;
		min-height: 60vh;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		background: linear-gradient(135deg, #0a0f1a 0%, #1a1f2e 50%, #0f172a 100%);
	}

	.hero-bg {
		position: absolute;
		inset: 0;
		background:
			radial-gradient(ellipse at 20% 50%, rgba(245,158,11,0.08) 0%, transparent 50%),
			radial-gradient(ellipse at 80% 50%, rgba(220,38,38,0.05) 0%, transparent 50%);
	}

	.hero-speed {
		position: absolute; inset: 0; opacity: 0.3;
		background: repeating-linear-gradient(
			90deg,
			rgba(245, 158, 11, 0.04) 0,
			rgba(245, 158, 11, 0.04) 2px,
			transparent 2px,
			transparent 20px
		);
		background-size: 40px 100%;
		animation: hero-speed-move 0.8s linear infinite;
	}

	@keyframes hero-speed-move {
		0% { background-position: 0 0; }
		100% { background-position: 40px 0; }
	}

	.hero-content {
		position: relative;
		z-index: 1;
		text-align: center;
		padding: 2rem 1rem;
		max-width: 800px;
	}

	.hero-logo { margin-bottom: 1.5rem; }

	.emblema {
		width: 100px; height: auto; border-radius: 50%;
		border: 3px solid var(--racing-amber);
		box-shadow: 0 0 20px rgba(245,158,11,0.3);
	}

	.hero-title { display: flex; flex-direction: column; gap: 0.25rem; margin-bottom: 1rem; }
	.title-line { display: block; }
	.title-line:first-child { font-size: 3rem; color: var(--racing-amber); text-shadow: 0 0 20px rgba(245,158,11,0.3); }
	.subtitle { font-size: 1.25rem; color: var(--racing-text-dim); font-family: 'Inter', sans-serif; text-transform: none; letter-spacing: 0.2em; }
	.hero-desc { color: var(--racing-text-dim); font-size: 1.1rem; margin-bottom: 2rem; }

	.hero-actions { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; }

	.btn-racing-primary {
		display: inline-block; padding: 0.75rem 2rem;
		background: linear-gradient(135deg, var(--racing-amber), var(--racing-amber-light));
		color: var(--racing-black); font-weight: 800; font-size: 1.1rem;
		text-decoration: none; text-transform: uppercase; letter-spacing: 0.05em;
		border-radius: 0.25rem; box-shadow: 0 0 20px rgba(245,158,11,0.3);
		transition: all 0.2s;
	}
	.btn-racing-primary:hover { transform: translateY(-2px); box-shadow: 0 0 30px rgba(245,158,11,0.5); }

	.hero-stripe {
		position: absolute; bottom: 0; left: 0; right: 0; height: 4px;
		background: repeating-linear-gradient(90deg, var(--racing-amber) 0, var(--racing-amber) 10px, var(--racing-black) 10px, var(--racing-black) 20px);
	}

	/* ─── Section ─────────────────────────────── */
	.info-section { padding: 3rem 1rem; background: var(--racing-dark); }
	.cards-section { background: var(--racing-black); }

	.section-title { text-align: center; margin-bottom: 2rem; }
	.section-title h2 { color: var(--racing-amber); font-size: 1.5rem; margin: 0; }
	.section-desc { color: var(--racing-text-dim); font-size: 0.95rem; margin: 0.5rem 0 0; }

	/* Skeleton loader */
	.skeleton-card { height: 80px; background: var(--racing-gray); border-radius: 0.5rem; position: relative; overflow: hidden; }
	.skeleton-card::after {
		content: ''; position: absolute; inset: 0;
		background: repeating-linear-gradient(90deg, transparent 0, rgba(245,158,11,0.05) 50%, transparent 100%);
		background-size: 200% 100%;
		animation: skeleton-loading 1.5s ease-in-out infinite;
	}
	@keyframes skeleton-loading {
		0% { background-position: 200% 0; }
		100% { background-position: -200% 0; }
	}

	.empty-carreras { text-align: center; padding: 3rem 1rem; }
	.empty-img { width: 80px; border-radius: 50%; opacity: 0.4; margin-bottom: 1rem; }
	.empty-carreras p { color: var(--racing-text-dim); margin: 0; }
	.empty-sub { font-size: 0.85rem; margin-top: 0.5rem !important; }

	/* ─── Categorias grid ─────────────────────── */
	.categorias-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 1rem;
		max-width: 900px;
		margin: 0 auto;
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
	.cat-badge-row { display: flex; align-items: center; gap: 0.5rem; margin-top: 0.5rem; }
	.cat-cta { color: var(--racing-amber); font-size: 0.85rem; font-weight: 600; }

	.badge-live { font-size: 0.75rem; color: #fca5a5; background: rgba(220,38,38,0.2); padding: 0.15rem 0.5rem; border-radius: 1rem; font-weight: 700; }
	.badge-done { font-size: 0.75rem; color: #a3e635; background: rgba(101,163,13,0.2); padding: 0.15rem 0.5rem; border-radius: 1rem; font-weight: 700; }

	.categoria-card.live { border-color: #ef4444; }
	.categoria-card.finalizado { border-color: #65a30d; }

	/* ─── Info Cards ──────────────────────────── */
	.info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 2rem;
		max-width: 1000px;
		margin: 0 auto;
	}

	.info-card {
		display: flex; flex-direction: column; align-items: center;
		background: var(--racing-gray); border: 1px solid var(--racing-border);
		border-radius: 0.5rem; padding: 1.5rem; text-align: center;
		text-decoration: none; transition: all 0.3s;
	}

	.info-card:hover { border-color: var(--racing-amber); transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.3); }

	.card-link { margin-top: 0.75rem; color: var(--racing-amber); font-weight: 600; font-size: 0.85rem; }

	.card-img { width: 80px; height: 80px; object-fit: cover; border-radius: 50%; margin-bottom: 1rem; border: 2px solid var(--racing-border); }
	.info-card h3 { color: var(--racing-amber); margin: 0 0 0.5rem 0; font-size: 1.25rem; }
	.info-card p { color: var(--racing-text-dim); font-size: 0.9rem; margin: 0; line-height: 1.5; }

	/* ─── Sponsors ────────────────────────────── */
	.sponsors { padding: 3rem 1rem; background: var(--racing-black); }

	.sponsor-logos {
		display: flex; justify-content: center; align-items: center; gap: 2rem;
		flex-wrap: wrap; max-width: 800px; margin: 0 auto;
	}

	.sponsor-img { height: 50px; width: auto; opacity: 0.5; transition: opacity 0.2s; }
	.sponsor-img:hover { opacity: 0.8; }

	@media (max-width: 640px) {
		.title-line:first-child { font-size: 2rem; }
		.emblema { width: 80px; }
		.hero-desc { font-size: 0.95rem; }
		.categorias-grid { grid-template-columns: 1fr; }
	}
</style>

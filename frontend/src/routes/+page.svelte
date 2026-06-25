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

<!-- ════ HERO — Arcade Racing Scene ════ -->
<section class="hero">
	<div class="hero-bg">
		<!-- Cielo atardecer -->
		<div class="sky"></div>

		<!-- Montañas lejanas -->
		<div class="mountains"></div>

		<!-- Estación de tren -->
		<div class="station">
			<div class="station-body"></div>
			<div class="station-roof"></div>
			<div class="station-sign">
				<span class="sign-text">BERAZATEGUI</span>
			</div>
			<div class="station-platform"></div>
		</div>

		<!-- Árboles pixel art -->
		<div class="tree t1"><div class="trunk"></div><div class="foliage"></div></div>
		<div class="tree t2"><div class="trunk"></div><div class="foliage"></div></div>
		<div class="tree t3"><div class="trunk"></div><div class="foliage"></div></div>
		<div class="tree t4"><div class="trunk"></div><div class="foliage"></div></div>

		<!-- Ruta en perspectiva -->
		<div class="road">
			<div class="road-line l1"></div>
			<div class="road-line l2"></div>
		</div>

		<!-- Tren azul y blanco animado -->
		<div class="train">
			<div class="train-body">
				<div class="train-cabin">
					<div class="cabin-window"></div>
				</div>
				<div class="train-stripe"></div>
				<div class="train-wheels">
					<div class="wheel w1"></div>
					<div class="wheel w2"></div>
					<div class="wheel w3"></div>
				</div>
			</div>
		</div>

		<!-- Líneas de velocidad -->
		<div class="hero-speed"></div>
	</div>

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
	/* ═══════════════════════════════════════════════════
	   ARCADE RACING SCENE — Hero
	   ═══════════════════════════════════════════════════ */

	.hero {
		position: relative;
		min-height: 65vh;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
	}

	/* ─── Contenedor de la escena ─────────────── */
	.hero-bg {
		position: absolute;
		inset: 0;
		background: linear-gradient(180deg, #1a1a3e 0%, #2d1b4e 20%, #e85d3a 45%, #f09030 55%, #2a3a1a 56%, #1a2a10 100%);
	}

	/* Cielo estrellado (puntitos) */
	.sky {
		position: absolute; inset: 0;
		background-image:
			radial-gradient(1px 1px at 10% 5%, rgba(255,255,255,0.6) 100%, transparent),
			radial-gradient(1px 1px at 25% 12%, rgba(255,255,255,0.4) 100%, transparent),
			radial-gradient(1.5px 1.5px at 40% 3%, rgba(255,255,255,0.7) 100%, transparent),
			radial-gradient(1px 1px at 55% 8%, rgba(255,255,255,0.5) 100%, transparent),
			radial-gradient(1px 1px at 70% 4%, rgba(255,255,255,0.6) 100%, transparent),
			radial-gradient(1.5px 1.5px at 85% 10%, rgba(255,255,255,0.4) 100%, transparent),
			radial-gradient(1px 1px at 15% 18%, rgba(255,255,255,0.3) 100%, transparent),
			radial-gradient(1px 1px at 60% 15%, rgba(255,255,255,0.5) 100%, transparent),
			radial-gradient(1px 1px at 90% 6%, rgba(255,255,255,0.4) 100%, transparent),
			radial-gradient(1.5px 1.5px at 50% 2%, rgba(255,255,255,0.6) 100%, transparent);
		pointer-events: none;
	}

	/* Montañas */
	.mountains {
		position: absolute; bottom: 44%; left: 0; right: 0; height: 20%;
		background:
			/* Montaña lejana 1 */
			linear-gradient(135deg, transparent 33%, #2a1a3a 33%, #2a1a3a 38%, transparent 38%),
			linear-gradient(225deg, transparent 33%, #2a1a3a 33%, #2a1a3a 38%, transparent 38%),
			/* Montaña lejana 2 */
			linear-gradient(135deg, transparent 50%, #3a2a4a 50%, #3a2a4a 55%, transparent 55%),
			linear-gradient(225deg, transparent 50%, #3a2a4a 50%, #3a2a4a 55%, transparent 55%),
			/* Montaña cercana */
			linear-gradient(135deg, transparent 60%, #1a2a1a 60%, #1a2a1a 68%, transparent 68%),
			linear-gradient(225deg, transparent 60%, #1a2a1a 60%, #1a2a1a 68%, transparent 68%);
		background-size: 100% 100%;
		background-repeat: no-repeat;
		background-position: 0 0;
		pointer-events: none;
	}

	/* ─── Árboles pixel art ──────────────────── */
	.tree {
		position: absolute; z-index: 2; pointer-events: none;
	}
	.trunk {
		position: absolute; bottom: 0; left: 50%; transform: translateX(-50%);
		width: 6px; height: 20px; background: #4a3520;
	}
	.foliage {
		position: absolute; bottom: 18px; left: 50%; transform: translateX(-50%);
		width: 24px; height: 24px;
		background: #2d5a1e;
		box-shadow:
			4px 4px 0 #3a6e28,
			-4px 4px 0 #3a6e28,
			0 -4px 0 #3a6e28,
			4px 0 0 #3a6e28,
			-4px 0 0 #3a6e28,
			0 4px 0 #3a6e28,
			8px 8px 0 #4a8030,
			-8px 8px 0 #4a8030,
			0 -8px 0 #4a8030;
	}

	.t1 { bottom: 42%; left: 5%; }
	.t2 { bottom: 40%; left: 18%; }
	.t3 { bottom: 38%; right: 22%; }
	.t4 { bottom: 41%; right: 6%; }
	.t1 .foliage { width: 28px; height: 28px; }
	.t1 .trunk { height: 26px; }

	/* ─── Estación de tren ────────────────────── */
	.station {
		position: absolute;
		bottom: 42%;
		right: 8%;
		width: 160px;
		height: 110px;
		z-index: 3;
		pointer-events: none;
	}

	/* Pared principal */
	.station-body {
		position: absolute;
		bottom: 16px; left: 15px; right: 15px; height: 65px;
		background: #c49a6c;
		box-shadow:
			-5px 0 0 #8a6a3a,
			5px 0 0 #8a6a3a,
			0 -3px 0 #d4aa7c;
	}

	/* Pilares de la estación */
	.station-body::before {
		content: '';
		position: absolute; top: 0; left: 10px; right: 10px; bottom: 0;
		background:
			/* Pilar izquierdo */
			linear-gradient(90deg, #8a6a3a 0, #8a6a3a 6px, transparent 6px),
			/* Pilar derecho */
			linear-gradient(90deg, transparent calc(100% - 6px), #8a6a3a calc(100% - 6px), #8a6a3a 100%);
	}

	/* Ventanas iluminadas */
	.station-body::after {
		content: '';
		position: absolute; top: 14px; left: 20px; right: 20px;
		height: 18px;
		background:
			linear-gradient(90deg,
				transparent 0, transparent 6px,
				#ffd700 6px, #ffd700 22px,
				transparent 22px, transparent 28px,
				#ffd700 28px, #ffd700 44px,
				transparent 44px, transparent 50px,
				#ffd700 50px, #ffd700 66px,
				transparent 66px, transparent 72px,
				#ffd700 72px, #ffd700 88px,
				transparent 88px);
		opacity: 0.7;
		box-shadow: 0 0 8px rgba(255,215,0,0.3);
	}

	/* Techo */
	.station-roof {
		position: absolute;
		top: 0; left: 0; right: 0; height: 20px;
		background: #5a3a1a;
		clip-path: polygon(3% 100%, 0% 10%, 100% 10%, 97% 100%);
	}

	/* Andén */
	.station-platform {
		position: absolute;
		bottom: 0; left: -15px; right: -15px; height: 8px;
		background: #777;
		box-shadow: 0 3px 0 #999, 0 5px 0 #555;
	}

	/* ─── Cartel BERAZATEGUI ──────────────────── */
	.station-sign {
		position: absolute;
		top: -28px; left: 50%; transform: translateX(-50%);
		background: #111;
		padding: 5px 14px;
		z-index: 5;
		box-shadow: 4px 4px 0 rgba(0,0,0,0.6);
		white-space: nowrap;
		border: 2px solid #444;
	}

	.sign-text {
		font-family: 'VT323', monospace;
		font-size: 1rem;
		color: #ffffff;
		font-weight: 700;
		letter-spacing: 0.2em;
	}

	/* ─── Tren azul y blanco ──────────────────── */
	.train {
		position: absolute;
		bottom: 48%;
		right: -200px;
		width: 120px;
		height: 40px;
		z-index: 4;
		pointer-events: none;
		animation: train-move 12s linear infinite;
	}

	@keyframes train-move {
		0% { transform: translateX(0); }
		100% { transform: translateX(-120vw); }
	}

	.train-body {
		position: absolute;
		bottom: 12px; left: 0; right: 40px; height: 28px;
		background: #1a5a9e;
		border-radius: 3px 0 0 3px;
		box-shadow: 0 -2px 0 #0d3a6e;
	}

	/* Cabina del conductor */
	.train-cabin {
		position: absolute;
		right: -36px; top: -8px;
		width: 36px; height: 36px;
		background: #1a5a9e;
		border-radius: 0 4px 4px 0;
		box-shadow: 0 -2px 0 #0d3a6e, 2px 0 0 #0d3a6e;
	}

	/* Ventana de cabina */
	.cabin-window {
		position: absolute;
		top: 8px; left: 6px; right: 4px; height: 14px;
		background: #a0d4ff;
		border-radius: 2px;
		box-shadow: inset 0 0 4px rgba(0,0,0,0.2);
	}

	/* Franja blanca decorativa */
	.train-stripe {
		position: absolute;
		bottom: 6px; left: 2px; right: 38px;
		height: 4px;
		background: #ffffff;
	}

	/* Ruedas */
	.train-wheels {
		position: absolute;
		bottom: 0; left: 8px; right: 34px;
		display: flex;
		justify-content: space-around;
	}

	.wheel {
		width: 10px; height: 10px;
		background: #333;
		border-radius: 50%;
		border: 2px solid #666;
		box-shadow: inset 0 0 2px rgba(0,0,0,0.5);
		animation: wheel-spin 1s linear infinite;
	}

	@keyframes wheel-spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	/* ─── Ruta en perspectiva ─────────────────── */
	.road {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 70%;
		height: 44%;
		background: #3a3a3a;
		clip-path: polygon(30% 0%, 70% 0%, 85% 100%, 15% 100%);
		z-index: 1;
		pointer-events: none;
	}

	/* Líneas de la ruta (divisoria central) */
	.road-line {
		position: absolute;
		left: 50%; transform: translateX(-50%);
		width: 4px;
		background: repeating-linear-gradient(180deg, var(--orange) 0px, var(--orange) 10px, transparent 10px, transparent 20px);
		opacity: 0.7;
	}
	.road-line.l1 { bottom: 0; height: 70%; }
	.road-line.l2 { bottom: 65%; height: 35%; width: 3px; opacity: 0.4; }

	/* ─── Speed lines ─────────────────────────── */
	.hero-speed {
		position: absolute; inset: 0; z-index: 5;
		background: repeating-linear-gradient(
			90deg,
			rgba(255, 255, 255, 0.03) 0,
			rgba(255, 255, 255, 0.03) 1px,
			transparent 1px,
			transparent 30px
		);
		background-size: 60px 100%;
		animation: hero-speed-move 0.6s linear infinite;
		pointer-events: none;
	}

	@keyframes hero-speed-move {
		0% { background-position: 0 0; }
		100% { background-position: 60px 0; }
	}

	/* ─── Hero content ────────────────────────── */
	.hero-content {
		position: relative;
		z-index: 10;
		text-align: center;
		padding: 2rem 1rem;
		max-width: 800px;
	}

	.hero-logo { margin-bottom: 1.5rem; }

	.emblema {
		width: 100px; height: auto; border-radius: 50%;
		border: 3px solid var(--orange);
		box-shadow: 0 0 20px rgba(245,158,11,0.5);
	}

	.hero-title { display: flex; flex-direction: column; gap: 0.25rem; margin-bottom: 1rem; }
	.title-line { display: block; }
	.title-line:first-child {
		font-size: 3rem;
		color: #ffffff;
		text-shadow:
			0 0 20px rgba(245,158,11,0.6),
			0 0 40px rgba(245,158,11,0.3),
			2px 2px 0 rgba(0,0,0,0.5),
			-1px -1px 0 rgba(0,0,0,0.3);
	}
	.subtitle {
		font-size: 1.25rem;
		color: #ffffff;
		font-family: 'Inter', sans-serif;
		text-transform: none;
		letter-spacing: 0.2em;
		text-shadow: 1px 1px 0 rgba(0,0,0,0.6);
	}
	.hero-desc {
		color: #f0f0f0;
		font-size: 1.1rem;
		margin-bottom: 2rem;
		text-shadow: 1px 1px 0 rgba(0,0,0,0.5);
	}

	.hero-actions { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; }

	.btn-racing-primary {
		display: inline-block; padding: 0.75rem 2rem;
		background: linear-gradient(135deg, var(--orange), var(--orange-glow));
		color: var(--arcade-black); font-weight: 800; font-size: 1.1rem;
		text-decoration: none; text-transform: uppercase; letter-spacing: 0.05em;
		border: 2px solid var(--orange-border);
		box-shadow: 3px 3px 0 0 var(--orange-border), 0 0 20px rgba(245,158,11,0.3);
		transition: all 0.2s;
	}
	.btn-racing-primary:hover { transform: translateY(-3px); box-shadow: 3px 3px 0 0 var(--orange), 0 0 30px rgba(245,158,11,0.5); }

	.hero-stripe {
		position: absolute; bottom: 0; left: 0; right: 0; height: 4px;
		background: repeating-linear-gradient(90deg, var(--orange) 0, var(--orange) 10px, transparent 10px, transparent 20px);
		z-index: 10;
	}

	/* ─── Section ─────────────────────────────── */
	.info-section { padding: 3rem 1rem; background: var(--arcade-dark); }
	.cards-section { background: var(--arcade-black); }

	.section-title { text-align: center; margin-bottom: 2rem; }
	.section-title h2 { color: var(--orange); font-size: 1.5rem; margin: 0; }
	.section-desc { color: var(--text-dim); font-size: 0.95rem; margin: 0.5rem 0 0; }

	/* Skeleton loader */
	.skeleton-card { height: 80px; background: var(--arcade-surface); border-radius: 0.5rem; position: relative; overflow: hidden; }
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
	.empty-carreras p { color: var(--text-dim); margin: 0; }
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
		background: var(--arcade-surface); border: 1px solid var(--border-color);
		border-radius: 0.5rem; padding: 1.5rem;
		text-decoration: none; transition: all 0.2s;
	}

	.categoria-card:hover {
		border-color: var(--orange);
		transform: translateY(-3px);
		box-shadow: 0 8px 24px rgba(0,0,0,0.3);
	}

	.cat-nombre { color: var(--orange); font-weight: 700; font-size: 1.1rem; }
	.cat-edades { color: var(--text-dim); font-size: 0.85rem; }
	.cat-badge-row { display: flex; align-items: center; gap: 0.5rem; margin-top: 0.5rem; }
	.cat-cta { color: var(--orange); font-size: 0.85rem; font-weight: 600; }

	.badge-live { font-size: 0.75rem; color: #fca5a5; background: rgba(220,38,38,0.2); padding: 0.15rem 0.5rem; border-radius: 1rem; font-weight: 700; }
	.badge-done { font-size: 0.75rem; color: #a3e635; background: rgba(101,163,13,0.2); padding: 0.15rem 0.5rem; border-radius: 1rem; font-weight: 700; }

	.categoria-card.live { border-color: var(--red-race-light); }
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
		background: var(--arcade-surface); border: 1px solid var(--border-color);
		border-radius: 0.5rem; padding: 1.5rem; text-align: center;
		text-decoration: none; transition: all 0.3s;
	}

	.info-card:hover { border-color: var(--orange); transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.3); }

	.card-link { margin-top: 0.75rem; color: var(--orange); font-weight: 600; font-size: 0.85rem; }

	.card-img { width: 80px; height: 80px; object-fit: cover; border-radius: 50%; margin-bottom: 1rem; border: 2px solid var(--border-color); }
	.info-card h3 { color: var(--orange); margin: 0 0 0.5rem 0; font-size: 1.25rem; }
	.info-card p { color: var(--text-dim); font-size: 0.9rem; margin: 0; line-height: 1.5; }

	/* ─── Sponsors ────────────────────────────── */
	.sponsors { padding: 3rem 1rem; background: var(--arcade-black); }

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

<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { page } from '$app/stores';
	import Podium from '$lib/components/Podium.svelte';
	import Celebration from '$lib/components/Celebration.svelte';

	const categoriaId = $derived($page.params.id);

	let fixture = $state<any>(null);
	let posiciones = $state<any[]>([]);
	let categoria = $state<any>(null);
	let autoNombres = $state<Record<string, string>>({});
	let autoNumeros = $state<Record<string, number>>({});
	let autosLista = $state<any[]>([]);
	let loading = $state(true);

	let podiumHeat = $state<{orden: string[]; label: string} | null>(null);
	let showCelebration = $state(false);
	let polling = $state(false);

	onMount(() => {
		cargarDatos();
		// Polling cada 10s si está en vivo
		const interval = setInterval(() => {
			if (fixture?.estado !== 'finalizado') {
				cargarDatosSilencioso();
			}
		}, 10000);
		return () => clearInterval(interval);
	});

	async function cargarDatosSilencioso() {
		polling = true;
		try {
			const [fixtureRes, posRes] = await Promise.all([
				fetch(`/api/categorias/${categoriaId}/fixture`),
				fetch(`/api/categorias/${categoriaId}/posiciones`)
			]);
			if (fixtureRes.ok) fixture = await fixtureRes.json();
			if (posRes.ok) posiciones = await posRes.json();
		} catch {} finally {
			polling = false;
		}
	}

	async function cargarDatos() {
		loading = true;
		try {
			const [fixtureRes, posRes, catRes, autosRes] = await Promise.all([
				fetch(`/api/categorias/${categoriaId}/fixture`),
				fetch(`/api/categorias/${categoriaId}/posiciones`),
				fetch(`/api/categorias/${categoriaId}`),
				fetch(`/api/categorias/${categoriaId}/autos`)
			]);

			if (fixtureRes.ok) fixture = await fixtureRes.json();
			if (posRes.ok) posiciones = await posRes.json();
			if (catRes.ok) categoria = await catRes.json();

			if (autosRes.ok) {
				const autos = await autosRes.json();
				autosLista = autos;
				for (const a of autos) {
					autoNombres[a.id] = a.nombre;
					autoNumeros[a.id] = a.numero;
				}
			}
		} catch (e) {
			// silent
		} finally {
			loading = false;
		}
	}

	let finalWinner = $derived.by(() => {
		if (!fixture?.heats) return null;
		const completados = fixture.heats.filter((h: any) => h.completado && h.orden_llegada?.length);
		if (completados.length === 0) return null;
		const last = completados[completados.length - 1];
		return {
			id: last.orden_llegada[0],
			nombre: autoNombres[last.orden_llegada[0]] || last.orden_llegada[0],
			numero: autoNumeros[last.orden_llegada[0]] || '?'
		};
	});

	// Top 3 del podio final
	let top3 = $derived.by(() => {
		if (!posiciones || posiciones.length === 0 || fixture?.estado !== 'finalizado') return [];
		return posiciones.slice(0, 3).map((p: any) => ({
			id: p.auto_id,
			nombre: autoNombres[p.auto_id] || p.nombre || p.auto_id,
			numero: autoNumeros[p.auto_id] ?? p.numero,
			puntos: p.puntos
		}));
	});

	function showPodium(heat: any) {
		podiumHeat = { orden: heat.orden_llegada, label: `Heat #${heat.numero}` };
	}
</script>

<svelte:head>
	<title>{categoria?.nombre ?? 'Carrera'} — Fixture D15</title>
</svelte:head>

<div class="carrera-page">
	<header class="carrera-header">
		<div class="header-inner">
			<a href="/carreras" class="back-link">← Categorías</a>
			<img src="/assets/derby_logo.jpg" alt="Derby" class="logo" />
			<div>
				<h1>{categoria?.nombre ?? 'Categoría'}</h1>
				<p class="subtitle">Pinewood Derby — Destacamento 15</p>
			</div>
		</div>
		<div class="header-stripe"></div>
	</header>

	{#if loading}
		<div class="skeleton-grid">
			<div class="skeleton-card"></div>
			<div class="skeleton-card"></div>
			<div class="skeleton-card"></div>
		</div>
	{/if}

	{#if fixture}
		<div class="status-bar" in:fade>
			<span class="status-badge" class:live={fixture.estado !== 'finalizado'}>
				{fixture.estado === 'finalizado' ? '🏁 Finalizado' : '🔴 En Vivo'}
				{#if polling}<span class="live-dot"></span>{/if}
			</span>
			<span class="heats-info">{fixture.heats?.length ?? 0} carreras</span>
		</div>
	{/if}

	{#if fixture?.estado === 'finalizado' && top3.length > 0}
		<section class="section podium-section" in:fly={{ y: 20, duration: 500 }}>
			<div class="podium-header">
				<h2>🏆 Podio Final</h2>
				<button class="btn-celebrate" onclick={() => showCelebration = true}>🎉 Celebrar</button>
			</div>
			<div class="podio">
				{#if top3[1]}
					<div class="podio-step second">
						<span class="podio-medal">🥈</span>
						<span class="podio-numero">#{top3[1].numero}</span>
						<span class="podio-nombre">{top3[1].nombre}</span>
						<span class="podio-puntos">{top3[1].puntos} pts</span>
						<div class="podio-bar" style="height: 80px">2°</div>
					</div>
				{/if}
				{#if top3[0]}
					<div class="podio-step first">
						<span class="podio-medal">🥇</span>
						<span class="podio-numero">#{top3[0].numero}</span>
						<span class="podio-nombre">{top3[0].nombre}</span>
						<span class="podio-puntos">{top3[0].puntos} pts</span>
						<div class="podio-bar" style="height: 120px">1°</div>
					</div>
				{/if}
				{#if top3[2]}
					<div class="podio-step third">
						<span class="podio-medal">🥉</span>
						<span class="podio-numero">#{top3[2].numero}</span>
						<span class="podio-nombre">{top3[2].nombre}</span>
						<span class="podio-puntos">{top3[2].puntos} pts</span>
						<div class="podio-bar" style="height: 60px">3°</div>
					</div>
				{/if}
			</div>
		</section>
	{/if}

	{#if posiciones?.length > 0}
		<section class="section" in:fade>
			<h2>📊 Tabla de Posiciones</h2>
			<div class="table-container">
				<table>
					<thead><tr>
						<th>#</th><th>Auto</th><th>N°</th><th>Pts</th><th>1°</th><th>2°</th><th>3°</th>
					</tr></thead>
					<tbody>
						{#each posiciones as pos, i}
							<tr class:top4={i < 4} class:champion={i === 0 && fixture?.estado === 'finalizado'}>
								<td class="pos">{i + 1}</td>
								<td>{autoNombres[pos.auto_id] || pos.auto_id}</td>
								<td>#{autoNumeros[pos.auto_id] ?? pos.numero}</td>
								<td class="pts">{pos.puntos}</td>
								<td>{pos.posiciones?.['1'] ?? 0}</td>
								<td>{pos.posiciones?.['2'] ?? 0}</td>
								<td>{pos.posiciones?.['3'] ?? 0}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</section>
	{/if}

	<!-- ─── Autos Registrados ─── -->
	{#if autosLista.length > 0}
		<section class="section" in:fade>
			<h2>🏎️ Autos Registrados</h2>
			<div class="autos-grid">
				{#each autosLista as auto (auto.id)}
					<div class="auto-card">
						{#if auto.foto_url}
							<img src={auto.foto_url} alt={auto.nombre} class="auto-foto" />
						{/if}
						<div class="auto-info">
							<span class="auto-numero">#{auto.numero}</span>
							<span class="auto-nombre">{auto.nombre}</span>
							<span class="auto-creador">por {auto.creador}</span>
							{#if auto.peso > 0}<span class="auto-peso">{auto.peso}g</span>{/if}
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if fixture?.heats?.length > 0}
		<section class="section" in:fade>
			<h2>🏎️ Carreras</h2>
			<div class="heats-grid">
				{#each fixture.heats as heat (heat.id)}
					<div class="heat-card" class:completed={heat.completado}>
						<div class="heat-header">
							<span class="heat-title">Carrera #{heat.numero}</span>
							{#if heat.completado}
								<button class="btn-podium" onclick={() => showPodium(heat)}>🥇</button>
							{:else}
								<span class="heat-pending">⏳ Pendiente</span>
							{/if}
						</div>
						<div class="heat-autos">
							{#each heat.auto_ids as autoId, i}
								<div class="heat-auto">
									<span class="auto-pos">{i + 1}.</span>
									<span class="auto-numero">#{autoNumeros[autoId] ?? '?'}</span>
									<span class="auto-nombre">{autoNombres[autoId] || autoId.slice(0, 8)}</span>
								</div>
							{/each}
						</div>
						{#if heat.completado && heat.orden_llegada?.length}
							<div class="heat-result">{heat.orden_llegada.map((id: string) => autoNumeros[id] ?? id).join(' → ')}</div>
						{/if}
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if !fixture && !loading}
		<section class="section" in:fade>
			<div class="empty-state">
				<img src="/assets/emblema.jpg" alt="Emblema" class="empty-img" />
				<h2>No hay carreras activas</h2>
				<p>Esperando que el administrador genere el fixture</p>
			</div>
		</section>
	{/if}

	<footer class="carrera-footer">
		<div class="footer-stripe"></div>
		<p>Destacamento 15 — Iglesia Betel</p>
	</footer>
</div>

<Podium ordenLlegada={podiumHeat?.orden ?? []} autoNombres={autoNombres} autoNumeros={autoNumeros} label={podiumHeat?.label ?? ''} show={podiumHeat !== null} onclose={() => podiumHeat = null} />
<Celebration winner={finalWinner?.nombre ?? ''} winnerNumero={finalWinner?.numero ?? 0} autoNombres={autoNombres} show={showCelebration} onclose={() => showCelebration = false} />

<style>
	.carrera-page { min-height: 100vh; background: var(--racing-black); }
	.carrera-header { background: linear-gradient(180deg, #0a0f1a, #0f172a); padding: 1rem; }
	.header-inner { display: flex; align-items: center; gap: 1rem; max-width: 800px; margin: 0 auto; flex-wrap: wrap; }
	.back-link { color: var(--racing-text-dim); text-decoration: none; font-size: 0.85rem; width: 100%; }
	.back-link:hover { color: var(--racing-amber); }
	.logo { height: 48px; width: auto; border-radius: 0.25rem; }
	h1 { color: var(--racing-amber); font-size: 1.5rem; margin: 0; }
	.subtitle { color: var(--racing-text-dim); font-size: 0.8rem; margin: 0; }
	.header-stripe { height: 3px; margin-top: 0.75rem; background: repeating-linear-gradient(90deg, var(--racing-amber) 0, var(--racing-amber) 8px, transparent 8px, transparent 16px); }

	.loading { text-align: center; color: #64748b; padding: 3rem; }

	/* Skeleton loader */
	.skeleton-grid { display: grid; gap: 1rem; max-width: 600px; margin: 2rem auto; padding: 0 1rem; }
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

	.live-dot {
		display: inline-block; width: 6px; height: 6px; border-radius: 50%;
		background: #22c55e; margin-left: 0.4rem;
		animation: pulse-dot 1.5s ease-in-out infinite;
		vertical-align: middle;
	}
	@keyframes pulse-dot {
		0%, 100% { opacity: 1; transform: scale(1); }
		50% { opacity: 0.4; transform: scale(0.8); }
	}

	.status-bar { display: flex; align-items: center; gap: 1rem; padding: 0.75rem 1rem; max-width: 800px; margin: 0 auto; }
	.status-badge { padding: 0.25rem 0.75rem; border-radius: 1rem; font-size: 0.8rem; font-weight: 700; background: #334155; color: #94a3b8; }
	.status-badge.live { background: rgba(220,38,38,0.2); color: #fca5a5; }
	.heats-info { color: var(--racing-text-dim); font-size: 0.85rem; }

	/* ─── Podio Final ─── */
	.podium-section { background: linear-gradient(180deg, #0f172a, #1a1f2e); border: 2px solid var(--racing-amber); border-radius: 0.75rem; padding: 2rem 1rem; }

	.podium-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
	.podium-header h2 { margin: 0; color: var(--racing-amber); font-size: 1.5rem; text-shadow: 0 0 20px rgba(245,158,11,0.3); }

	.btn-celebrate { padding: 0.5rem 1.25rem; background: linear-gradient(135deg, var(--racing-amber), #f59e0b); color: var(--racing-black); font-weight: 800; font-size: 0.9rem; border: none; border-radius: 0.25rem; cursor: pointer; transition: all 0.2s; }
	.btn-celebrate:hover { transform: scale(1.05); box-shadow: 0 0 20px rgba(245,158,11,0.4); }

	.podio { display: flex; justify-content: center; align-items: flex-end; gap: 1rem; max-width: 600px; margin: 0 auto; }

	.podio-step {
		display: flex; flex-direction: column; align-items: center; gap: 0.3rem;
		flex: 1; max-width: 160px;
	}

	.podio-medal { font-size: 2rem; }
	.podio-numero { color: var(--racing-amber); font-weight: 700; font-size: 1.2rem; }
	.podio-nombre { color: var(--racing-text); font-weight: 600; font-size: 0.9rem; text-align: center; }
	.podio-puntos { color: var(--racing-text-dim); font-size: 0.8rem; }

	.podio-bar {
		width: 100%; display: flex; align-items: center; justify-content: center;
		border-radius: 0.25rem 0.25rem 0 0; font-weight: 900; font-size: 1.2rem; color: var(--racing-black);
	}

	.first .podio-bar { background: linear-gradient(180deg, #ffd700, #f59e0b); }
	.second .podio-bar { background: linear-gradient(180deg, #c0c0c0, #94a3b8); }
	.third .podio-bar { background: linear-gradient(180deg, #cd7f32, #92400e); }

	.section { max-width: 800px; margin: 1.5rem auto; padding: 0 1rem; }
	.section h2 { color: var(--racing-amber); font-size: 1.2rem; margin-bottom: 1rem; }

	.table-container { overflow-x: auto; }
	table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
	th { background: var(--racing-dark); color: var(--racing-amber); padding: 0.75rem; text-align: left; font-size: 0.75rem; text-transform: uppercase; }
	td { padding: 0.6rem 0.75rem; border-top: 1px solid var(--racing-border); }
	.top4 td { background: rgba(5,150,105,0.05); }
	.champion td { background: rgba(245,158,11,0.15); font-weight: 700; }
	.pos { color: var(--racing-amber); font-weight: 700; }
	.pts { color: var(--racing-amber); font-weight: 700; font-size: 1rem; }

	.heats-grid { display: grid; gap: 0.75rem; }
	.heat-card { background: var(--racing-dark); border: 1px solid var(--racing-border); border-radius: 0.5rem; padding: 1rem; }
	.heat-card.completed { border-color: #059669; }
	.heat-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
	.heat-title { color: var(--racing-amber); font-weight: 700; }
	.heat-pending { color: var(--racing-text-dim); font-size: 0.8rem; }
	.btn-podium { background: none; border: none; font-size: 1.5rem; cursor: pointer; }
	.btn-podium:hover { transform: scale(1.2); }
	.heat-autos { display: flex; flex-direction: column; gap: 0.2rem; }
	.heat-auto { display: flex; gap: 0.5rem; font-size: 0.9rem; }
	.auto-pos { color: var(--racing-text-dim); width: 1.2rem; }
	.auto-numero { color: var(--racing-amber); font-weight: 700; min-width: 2rem; }
	.auto-nombre { color: var(--racing-text); }
	.heat-result { margin-top: 0.5rem; font-size: 0.8rem; color: #059669; padding: 0.5rem; background: rgba(5,150,105,0.1); border-radius: 0.25rem; }

	.empty-state { text-align: center; padding: 4rem 1rem; }
	.empty-img { width: 100px; border-radius: 50%; opacity: 0.5; margin-bottom: 1rem; }
	.empty-state h2 { color: var(--racing-text-dim); }
	.empty-state p { color: var(--racing-text-dim); font-size: 0.9rem; }

	/* ─── Autos grid ─── */
	.autos-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 0.75rem;
	}

	.auto-card {
		display: flex; align-items: center; gap: 0.75rem;
		background: var(--racing-dark); border: 1px solid var(--racing-border);
		border-radius: 0.5rem; padding: 0.75rem;
	}

	.auto-foto { width: 48px; height: 48px; object-fit: cover; border-radius: 0.25rem; border: 1px solid var(--racing-border); }

	.auto-info { display: flex; flex-direction: column; gap: 0.1rem; }
	.auto-numero { color: var(--racing-amber); font-weight: 700; font-size: 1rem; }
	.auto-nombre { color: var(--racing-text); font-size: 0.9rem; }
	.auto-creador { color: var(--racing-text-dim); font-size: 0.8rem; }
	.auto-peso { color: var(--racing-amber); font-size: 0.75rem; font-weight: 600; }

	.carrera-footer { text-align: center; padding: 2rem 1rem; color: var(--racing-text-dim); font-size: 0.8rem; }
	.footer-stripe { height: 2px; background: repeating-linear-gradient(90deg, var(--racing-red) 0, var(--racing-red) 6px, transparent 6px, transparent 12px); margin-bottom: 1rem; max-width: 800px; margin-left: auto; margin-right: auto; }

	@media (max-width: 640px) { h1 { font-size: 1.2rem; } .btn-champion { font-size: 1rem; padding: 0.75rem 1.5rem; } }
</style>

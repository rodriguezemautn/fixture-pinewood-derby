<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';

	const categoriaId = $derived($page.params.id);

	let fixture = $state<any>(null);
	let posiciones = $state<any[]>([]);
	let categoria = $state<any>(null);
	let autos = $state<any[]>([]);
	let autoMap = $state<Record<string, {nombre: string; numero: number}>>({});
	let loading = $state(true);
	let error = $state('');

	let rondas = $state(3);
	let actionLoading = $state(false);

	// Estado para registro de resultado
	let resultadoState = $state<Record<string, string[]>>({});

	onMount(() => {
		cargarDatos();
	});

	async function cargarDatos() {
		loading = true;
		error = '';
		try {
			const [fixtureRes, posRes, catRes, autosRes] = await Promise.all([
				fetch(`/api/categorias/${categoriaId}/fixture`),
				fetch(`/api/categorias/${categoriaId}/posiciones`),
				fetch(`/api/categorias/${categoriaId}`),
				fetch(`/api/categorias/${categoriaId}/autos`)
			]);
			if (fixtureRes.ok) {
				fixture = await fixtureRes.json();
				// Inicializar estado de resultado para heats sin completar
				if (fixture?.heats) {
					for (const h of fixture.heats) {
						if (!h.completado) {
							resultadoState[h.id] = [];
						}
					}
				}
			}
			if (posRes.ok) posiciones = await posRes.json();
			if (catRes.ok) categoria = await catRes.json();
			if (autosRes.ok) {
				autos = await autosRes.json();
				for (const a of autos) {
					autoMap[a.id] = { nombre: a.nombre, numero: a.numero };
				}
			}
		} catch (e) {
			error = 'Error al cargar datos. ¿El backend está corriendo?';
		} finally {
			loading = false;
		}
	}

	async function generarFixture() {
		actionLoading = true;
		error = '';
		try {
			const res = await apiFetch(`/api/categorias/${categoriaId}/fixture?rondas=${rondas}`, { method: 'POST' });
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al generar fixture';
				return;
			}
			await cargarDatos();
		} catch { error = 'Error de conexión'; }
		finally { actionLoading = false; }
	}

	function seleccionarPosicion(heatId: string, autoId: string, posicion: number) {
		const autosDisponibles = getAutosDisponibles(heatId, autoId);
		// Remover auto de su posición anterior si ya estaba asignado
		const estado = [...resultadoState[heatId]];
		const idxPrev = estado.indexOf(autoId);
		if (idxPrev !== -1) {
			estado.splice(idxPrev, 1);
		}
		// Insertar en la nueva posición
		const idx = posicion - 1;
		if (idx < estado.length) {
			estado[idx] = autoId;
		} else {
			// Llenar con vacíos hasta la posición
			while (estado.length < idx) estado.push('');
			estado.push(autoId);
		}
		resultadoState[heatId] = estado;
	}

	function getPosicionActual(heatId: string, autoId: string): number {
		return resultadoState[heatId]?.indexOf(autoId) + 1 || 0;
	}

	function getAutosDisponibles(heatId: string, exceptoId: string): string[] {
		const heat = fixture?.heats?.find((h: any) => h.id === heatId);
		if (!heat) return [];
		return heat.auto_ids.filter((id: string) =>
			id === exceptoId || !resultadoState[heatId]?.includes(id)
		);
	}

	function getAutoEnPosicion(heatId: string, pos: number): string | null {
		const estado = resultadoState[heatId];
		if (!estado || pos >= estado.length) return null;
		return estado[pos] || null;
	}

	async function registrarResultado(heatId: string) {
		const orden = resultadoState[heatId]?.filter(Boolean);
		if (!orden || orden.length < 2) {
			error = 'Asigná al menos 2 posiciones';
			return;
		}

		actionLoading = true;
		error = '';
		try {
			const res = await apiFetch(`/api/carreras/${heatId}/resultado`, {
				method: 'POST',
				body: JSON.stringify({ orden_llegada: orden })
			});
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al registrar';
				return;
			}
			await cargarDatos();
		} catch { error = 'Error de conexión'; }
		finally { actionLoading = false; }
	}

	async function generarFinal() {
		actionLoading = true;
		error = '';
		try {
			const res = await apiFetch(`/api/categorias/${categoriaId}/final`, { method: 'POST' });
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al generar final';
				return;
			}
			await cargarDatos();
		} catch { error = 'Error de conexión'; }
		finally { actionLoading = false; }
	}

	const puedeGenerarFinal = $derived(
		fixture?.heats?.length > 0 && fixture?.heats?.every((h: any) => h.completado)
	);

	function getAutoDisplay(autoId: string): string {
		const a = autoMap[autoId];
		return a ? `#${a.numero} ${a.nombre}` : autoId.slice(0, 8);
	}
</script>

<div class="page">
	<div class="header">
		<div>
			<a href="/admin/categorias/{categoriaId}/autos" class="back-link">← Autos</a>
			<h1>{categoria?.nombre ?? 'Cargando...'}</h1>
		</div>
		<a href="/carreras/{categoriaId}" class="btn btn-view" target="_blank">👁️ Vista Pública</a>
	</div>

	{#if error}
		<div class="error-banner" in:fade>{error}</div>
	{/if}

	{#if loading}
		<p class="loading">Cargando...</p>
	{/if}

	<!-- Generar Fixture -->
	{#if !loading && !fixture}
		<section class="card" in:fly={{ y: 10, duration: 300 }}>
			<h2>🏁 Generar Fixture</h2>
			<p class="desc">Generá la secuencia de carreras para determinar los 4 mejores autos.</p>
			<div class="form-row">
				<label>
					Rondas
					<input type="number" bind:value={rondas} min="1" max="10" class="input" />
				</label>
				<button class="btn btn-primary" onclick={generarFixture} disabled={actionLoading}>
					{actionLoading ? 'Generando...' : '🔥 Generar Fixture'}
				</button>
			</div>
		</section>
	{/if}

	<!-- Heats -->
	{#if fixture?.heats?.length > 0}
		<section class="card" in:fade>
			<div class="card-header">
				<h2>🏎️ Heats de Clasificación</h2>
				<span class="badge">{fixture.estado}</span>
			</div>

			<div class="heats-grid">
				{#each fixture.heats as heat (heat.id)}
					<div class="heat-card" class:completed={heat.completado} in:fade>
						<div class="heat-header">
							<span class="heat-num">Heat #{heat.numero}</span>
							{#if heat.completado}
								<span class="heat-status done">✅ Completado</span>
							{:else}
								<span class="heat-status pending">⏳ Pendiente</span>
							{/if}
						</div>

						<div class="heat-autos">
							{#each heat.auto_ids as autoId, i}
								<div class="heat-auto-row">
									<span class="heat-auto-label">{getAutoDisplay(autoId)}</span>
								</div>
							{/each}
						</div>

						{#if heat.completado && heat.orden_llegada?.length}
							<div class="heat-result">
								<strong>Podio:</strong>
								{#each heat.orden_llegada as id, i}
									<span class="podio-item">
										<span class="medal">{['🥇','🥈','🥉','4️⃣'][i]}</span>
										{getAutoDisplay(id)}
									</span>
									{#if i < heat.orden_llegada.length - 1}<span class="arrow">→</span>{/if}
								{/each}
							</div>
						{/if}

						{#if !heat.completado}
							<div class="result-form">
								<p class="result-hint">Asigná la posición de cada auto:</p>
								{#each heat.auto_ids as autoId}
									<div class="pos-row">
										<span class="pos-auto">{getAutoDisplay(autoId)}</span>
										<select
											class="pos-select"
											value={getPosicionActual(heat.id, autoId)}
											onchange={(e) => seleccionarPosicion(heat.id, autoId, parseInt(e.target.value))}
										>
											<option value={0}>—</option>
											{#each heat.auto_ids as _, pi}
												<option
													value={pi + 1}
													disabled={getAutoEnPosicion(heat.id, pi) !== null && getAutoEnPosicion(heat.id, pi) !== autoId}
												>
													{pi + 1}° lugar
												</option>
											{/each}
										</select>
									</div>
								{/each}
								<button class="btn btn-sm btn-primary" onclick={() => registrarResultado(heat.id)} disabled={actionLoading}>
									{actionLoading ? 'Guardando...' : '✅ Registrar Resultado'}
								</button>
							</div>
						{/if}
					</div>
				{/each}
			</div>

			{#if puedeGenerarFinal}
				<div class="final-section" in:fade>
					<button class="btn btn-racing" onclick={generarFinal} disabled={actionLoading}>🏆 Generar Carrera Final</button>
				</div>
			{/if}
		</section>
	{/if}

	<!-- Tabla de Posiciones -->
	{#if posiciones?.length > 0}
		<section class="card" in:fade>
			<h2>📊 Tabla de Posiciones</h2>
			<div class="table-container">
				<table>
					<thead><tr>
						<th>#</th><th>Auto</th><th>N°</th><th>Puntos</th><th>1°</th><th>2°</th><th>3°</th><th>4°</th><th>Carreras</th>
					</tr></thead>
					<tbody>
						{#each posiciones as pos, i}
							<tr class:top4={i < 4} class:winner={i === 0}>
								<td class="pos">{i + 1}</td>
								<td>{pos.nombre}</td>
								<td>#{pos.numero}</td>
								<td class="pts">{pos.puntos}</td>
								<td>{pos.posiciones?.['1'] ?? 0}</td>
								<td>{pos.posiciones?.['2'] ?? 0}</td>
								<td>{pos.posiciones?.['3'] ?? 0}</td>
								<td>{pos.posiciones?.['4'] ?? 0}</td>
								<td>{pos.carreras}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</section>
	{/if}
</div>

<style>
	.page { max-width: 1000px; margin: 0 auto; }
	.header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 2rem; }
	h1 { color: var(--racing-amber); font-size: 1.75rem; margin: 0.5rem 0 0 0; }
	.back-link { color: #64748b; text-decoration: none; font-size: 0.875rem; }
	.back-link:hover { color: var(--racing-amber); }
	.loading { text-align: center; color: #64748b; padding: 3rem; }

	.card { background: var(--racing-gray); border: 1px solid var(--racing-border); border-radius: 0.5rem; padding: 1.5rem; margin-bottom: 1.5rem; }
	.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
	h2 { color: var(--racing-amber); margin: 0 0 1rem 0; font-size: 1.25rem; }
	.card h2 { margin: 0 0 1rem 0; }
	.desc { color: var(--racing-text-dim); margin-bottom: 1rem; }
	.form-row { display: flex; gap: 1rem; align-items: flex-end; flex-wrap: wrap; }
	label { display: flex; flex-direction: column; gap: 0.25rem; color: #94a3b8; font-size: 0.875rem; }
	.input { padding: 0.5rem; border: 1px solid var(--racing-border); border-radius: 0.25rem; background: var(--racing-black); color: var(--racing-text); font-size: 1rem; }
	.input:focus { outline: none; border-color: var(--racing-amber); }

	.btn { padding: 0.5rem 1rem; border: none; border-radius: 0.25rem; cursor: pointer; font-weight: 600; font-size: 0.875rem; transition: all 0.2s; }
	.btn-primary { background: var(--racing-amber); color: var(--racing-black); }
	.btn-primary:hover { background: var(--racing-amber-light); }
	.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
	.btn-sm { padding: 0.35rem 0.75rem; font-size: 0.8rem; background: #334155; color: #e2e8f0; }
	.btn-sm:hover { background: #475569; }
	.btn-view { text-decoration: none; padding: 0.5rem 1rem; background: #334155; color: #e2e8f0; border-radius: 0.25rem; font-size: 0.875rem; }
	.btn-racing { background: linear-gradient(135deg, var(--racing-amber), var(--racing-amber-light)); color: var(--racing-black); font-weight: 800; font-size: 1.1rem; padding: 0.75rem 2rem; border: none; border-radius: 0.25rem; cursor: pointer; }
	.btn-racing:hover { transform: translateY(-2px); }
	.badge { padding: 0.25rem 0.75rem; border-radius: 1rem; font-size: 0.75rem; text-transform: uppercase; background: var(--racing-amber); color: var(--racing-black); font-weight: 700; }
	.error-banner { background: rgba(220,38,38,0.2); border: 1px solid var(--racing-red); color: #fca5a5; padding: 0.75rem 1rem; border-radius: 0.25rem; margin-bottom: 1rem; }

	/* Heats */
	.heats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1rem; }
	.heat-card { background: var(--racing-black); border: 1px solid var(--racing-border); border-radius: 0.5rem; padding: 1rem; transition: all 0.2s; }
	.heat-card:hover { border-color: var(--racing-amber); }
	.heat-card.completed { border-color: #059669; }
	.heat-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
	.heat-num { font-weight: 700; color: var(--racing-amber); }
	.heat-status { font-size: 0.8rem; }

	.heat-autos { display: flex; flex-direction: column; gap: 0.3rem; margin-bottom: 0.75rem; }
	.heat-auto-label { font-size: 0.9rem; color: var(--racing-text); padding: 0.2rem 0; }

	.heat-result { font-size: 0.85rem; color: #059669; margin-bottom: 0.75rem; padding: 0.75rem; background: rgba(5,150,105,0.1); border-radius: 0.25rem; display: flex; flex-wrap: wrap; align-items: center; gap: 0.3rem; }
	.podio-item { display: inline-flex; align-items: center; gap: 0.2rem; }
	.medal { font-size: 1.1rem; }
	.arrow { color: var(--racing-text-dim); margin: 0 0.2rem; }

	/* Result form */
	.result-form { margin-top: 0.75rem; padding-top: 0.75rem; border-top: 1px solid var(--racing-border); }
	.result-hint { font-size: 0.8rem; color: var(--racing-text-dim); margin: 0 0 0.5rem 0; }
	.pos-row { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.4rem; }
	.pos-auto { flex: 1; font-size: 0.9rem; color: var(--racing-text); }
	.pos-select {
		padding: 0.3rem 0.5rem; border: 1px solid var(--racing-border);
		border-radius: 0.25rem; background: var(--racing-dark);
		color: var(--racing-text); font-size: 0.85rem; min-width: 100px;
	}
	.pos-select:focus { outline: none; border-color: var(--racing-amber); }

	.final-section { margin-top: 1.5rem; text-align: center; padding-top: 1.5rem; border-top: 1px solid var(--racing-border); }

	/* Standings table */
	.table-container { overflow-x: auto; }
	table { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
	th { background: var(--racing-black); color: var(--racing-amber); padding: 0.75rem; text-align: left; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; }
	td { padding: 0.75rem; border-top: 1px solid var(--racing-border); }
	.top4 td { background: rgba(5,150,105,0.05); }
	.winner td { background: rgba(245,158,11,0.1); font-weight: 700; }
	.pos { color: var(--racing-amber); font-weight: 700; }
	.pts { color: var(--racing-amber); font-weight: 700; font-size: 1.1rem; }

	@media (max-width: 640px) { .heats-grid { grid-template-columns: 1fr; } }
</style>

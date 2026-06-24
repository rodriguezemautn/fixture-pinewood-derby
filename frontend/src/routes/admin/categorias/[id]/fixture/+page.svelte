<script lang="ts">
	import { fade, fly } from 'svelte/transition';

	let { data } = $props();
	let { fixture, posiciones, categoriaId, categoria } = $derived(data);

	let rondas = $state(3);
	let loading = $state(false);
	let error = $state('');

	let activeHeat = $state<string | null>(null);
	let ordenInput = $state('');

	async function generarFixture() {
		loading = true;
		error = '';
		try {
			const res = await fetch(`/api/categorias/${categoriaId}/fixture?rondas=${rondas}`, { method: 'POST' });
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al generar fixture';
				return;
			}
			window.location.reload();
		} catch { error = 'Error de conexión'; }
		finally { loading = false; }
	}

	async function registrarResultado(heatId: string) {
		const orden = ordenInput.split(',').map(s => s.trim()).filter(Boolean);
		if (orden.length < 2) return;

		loading = true;
		error = '';
		try {
			const res = await fetch(`/api/carreras/${heatId}/resultado`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ orden_llegada: orden })
			});
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al registrar';
				return;
			}
			activeHeat = null;
			ordenInput = '';
			window.location.reload();
		} catch { error = 'Error de conexión'; }
		finally { loading = false; }
	}

	async function generarFinal() {
		loading = true;
		error = '';
		try {
			const res = await fetch(`/api/categorias/${categoriaId}/final`, { method: 'POST' });
			if (!res.ok) {
				const d = await res.json();
				error = d.error || 'Error al generar final';
				return;
			}
			window.location.reload();
		} catch { error = 'Error de conexión'; }
		finally { loading = false; }
	}

	const puedeGenerarFinal = $derived(
		fixture?.heats?.length > 0 && fixture?.heats?.every((h: any) => h.completado)
	);
</script>

<div class="page">
	<div class="header">
		<div>
			<a href="/admin/categorias/{categoriaId}/autos" class="back-link">← Autos</a>
			<h1>{categoria?.nombre ?? 'Cargando...'}</h1>
		</div>
	</div>

	{#if error}
		<div class="error-banner" in:fade>{error}</div>
	{/if}

	<!-- ─── Generar Fixture ─── -->
	{#if !fixture}
		<section class="card" in:fly={{ y: 10, duration: 300 }}>
			<h2>🏁 Generar Fixture</h2>
			<p class="desc">Generá la secuencia de carreras para determinar los 4 mejores autos.</p>
			<div class="form-row">
				<label>
					Rondas
					<input type="number" bind:value={rondas} min="1" max="10" class="input" />
				</label>
				<button class="btn btn-primary" onclick={generarFixture} disabled={loading}>
					{loading ? 'Generando...' : '🔥 Generar Fixture'}
				</button>
			</div>
		</section>
	{/if}

	<!-- ─── Heats ─── -->
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
								<span class="heat-status done">✅</span>
							{:else}
								<span class="heat-status pending">⏳</span>
							{/if}
						</div>

						<div class="heat-autos">
							{#each heat.auto_ids as autoId, i}
								<div class="heat-auto">
									<span class="heat-auto-pos">{i + 1}.</span>
									<span class="heat-auto-id">{autoId.slice(0, 8)}</span>
								</div>
							{/each}
						</div>

						{#if heat.completado && heat.orden_llegada?.length}
							<div class="heat-result">
								<strong>Orden:</strong> {heat.orden_llegada.join(' → ')}
							</div>
						{/if}

						{#if !heat.completado}
							{#if activeHeat === heat.id}
								<div class="result-form">
									<input
										type="text"
										bind:value={ordenInput}
										placeholder="auto-id-1, auto-id-2, auto-id-3, auto-id-4"
										class="input"
									/>
									<div class="btn-row">
										<button class="btn btn-sm" onclick={() => { activeHeat = null; ordenInput = ''; }}>Cancelar</button>
										<button class="btn btn-sm btn-primary" onclick={() => registrarResultado(heat.id)} disabled={loading}>
											Registrar
										</button>
									</div>
								</div>
							{:else}
								<button class="btn btn-sm btn-primary" onclick={() => activeHeat = heat.id}>
									Registrar Resultado
								</button>
							{/if}
						{/if}
					</div>
				{/each}
			</div>

			{#if puedeGenerarFinal}
				<div class="final-section" in:fade>
					<button class="btn btn-racing" onclick={generarFinal} disabled={loading}>
						🏆 Generar Carrera Final
					</button>
				</div>
			{/if}
		</section>
	{/if}

	<!-- ─── Tabla de Posiciones ─── -->
	{#if posiciones?.length > 0}
		<section class="card" in:fade>
			<h2>📊 Tabla de Posiciones</h2>
			<div class="table-container">
				<table>
					<thead>
						<tr>
							<th>#</th>
							<th>Auto</th>
							<th>N°</th>
							<th>Puntos</th>
							<th>1°</th>
							<th>2°</th>
							<th>3°</th>
							<th>4°</th>
							<th>Carreras</th>
						</tr>
					</thead>
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

	.header { margin-bottom: 2rem; }
	h1 { color: var(--racing-amber); font-size: 1.75rem; margin: 0.5rem 0 0 0; }
	.back-link { color: #64748b; text-decoration: none; font-size: 0.875rem; }
	.back-link:hover { color: var(--racing-amber); }

	.card {
		background: var(--racing-gray);
		border: 1px solid var(--racing-border);
		border-radius: 0.5rem;
		padding: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	h2 { color: var(--racing-amber); margin: 0 0 1rem 0; font-size: 1.25rem; }
	.card h2 { margin: 0 0 1rem 0; }

	.desc { color: var(--racing-text-dim); margin-bottom: 1rem; }

	.form-row {
		display: flex;
		gap: 1rem;
		align-items: flex-end;
		flex-wrap: wrap;
	}

	label { display: flex; flex-direction: column; gap: 0.25rem; color: #94a3b8; font-size: 0.875rem; }
	.input {
		padding: 0.5rem; border: 1px solid var(--racing-border);
		border-radius: 0.25rem; background: var(--racing-black);
		color: var(--racing-text); font-size: 1rem;
	}
	.input:focus { outline: none; border-color: var(--racing-amber); }

	.btn {
		padding: 0.5rem 1rem; border: none; border-radius: 0.25rem;
		cursor: pointer; font-weight: 600; font-size: 0.875rem;
		transition: all 0.2s;
	}
	.btn-primary { background: var(--racing-amber); color: var(--racing-black); }
	.btn-primary:hover { background: var(--racing-amber-light); }
	.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
	.btn-sm { padding: 0.35rem 0.75rem; font-size: 0.8rem; background: #334155; color: #e2e8f0; }
	.btn-sm:hover { background: #475569; }

	.btn-racing {
		background: linear-gradient(135deg, var(--racing-amber), var(--racing-amber-light));
		color: var(--racing-black); font-weight: 800; font-size: 1.1rem;
		padding: 0.75rem 2rem;
	}
	.btn-racing:hover { transform: translateY(-2px); }

	.badge {
		padding: 0.25rem 0.75rem; border-radius: 1rem;
		font-size: 0.75rem; text-transform: uppercase;
		background: var(--racing-amber); color: var(--racing-black);
		font-weight: 700;
	}

	.error-banner {
		background: rgba(220,38,38,0.2); border: 1px solid var(--racing-red);
		color: #fca5a5; padding: 0.75rem 1rem; border-radius: 0.25rem;
		margin-bottom: 1rem;
	}

	/* Heats */
	.heats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }

	.heat-card {
		background: var(--racing-black); border: 1px solid var(--racing-border);
		border-radius: 0.5rem; padding: 1rem;
		transition: all 0.2s;
	}
	.heat-card:hover { border-color: var(--racing-amber); }
	.heat-card.completed { border-color: #059669; }

	.heat-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
	.heat-num { font-weight: 700; color: var(--racing-amber); }

	.heat-autos { display: flex; flex-direction: column; gap: 0.25rem; margin-bottom: 0.75rem; }
	.heat-auto { display: flex; gap: 0.5rem; font-size: 0.85rem; color: var(--racing-text-dim); }
	.heat-auto-pos { color: var(--racing-amber); font-weight: 600; width: 1.5rem; }

	.heat-result { font-size: 0.8rem; color: #059669; margin-bottom: 0.75rem; padding: 0.5rem; background: rgba(5,150,105,0.1); border-radius: 0.25rem; }

	.result-form { display: flex; flex-direction: column; gap: 0.5rem; }
	.btn-row { display: flex; gap: 0.5rem; }

	.final-section { margin-top: 1.5rem; text-align: center; padding-top: 1.5rem; border-top: 1px solid var(--racing-border); }

	/* Standings table */
	.table-container { overflow-x: auto; }
	table { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
	th {
		background: var(--racing-black); color: var(--racing-amber);
		padding: 0.75rem; text-align: left; font-size: 0.8rem;
		text-transform: uppercase; letter-spacing: 0.05em;
	}
	td { padding: 0.75rem; border-top: 1px solid var(--racing-border); }

	.top4 td { background: rgba(5,150,105,0.05); }
	.winner td { background: rgba(245,158,11,0.1); font-weight: 700; }
	.pos { color: var(--racing-amber); font-weight: 700; }
	.pts { color: var(--racing-amber); font-weight: 700; font-size: 1.1rem; }

	@media (max-width: 640px) {
		.heats-grid { grid-template-columns: 1fr; }
	}
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { page } from '$app/stores';
	import AutoForm from '$lib/components/AutoForm.svelte';
	import { apiFetch } from '$lib/api';

	const categoriaId = $derived($page.params.id);

	let autos = $state<any[]>([]);
	let categoria = $state<any>(null);
	let loading = $state(true);
	let error = $state('');
	let success = $state('');
	let competencias = $state<any[]>([]);
	let compLoading = $state(true);
	let genLoading = $state(false);
	let rondas = $state(3);

	let showForm = $state(false);
	let editingId = $state<string | null>(null);

	onMount(() => {
		cargarDatos();
	});

	async function cargarDatos() {
		loading = true;
		compLoading = true;
		error = '';
		try {
			const [autosRes, catRes, compRes] = await Promise.all([
				fetch(`/api/categorias/${categoriaId}/autos`),
				fetch(`/api/categorias/${categoriaId}`),
				fetch(`/api/categorias/${categoriaId}/competencias`)
			]);
			if (autosRes.ok) { autos = await autosRes.json(); }
			if (catRes.ok) { categoria = await catRes.json(); }
			if (compRes.ok) { competencias = await compRes.json(); compLoading = false; }
			else { compLoading = false; }
		} catch (e) {
			error = 'Error al cargar datos. ¿El backend está corriendo?';
		} finally {
			loading = false;
		}
	}

	function handleEdit(id: string) {
		editingId = id;
		showForm = true;
	}

	function handleClose() {
		showForm = false;
		editingId = null;
	}

	async function handleSaved() {
		showForm = false;
		editingId = null;
		await cargarDatos();
	}

	async function handleDelete(id: string) {
		if (!confirm('¿Eliminar este auto?')) return;
		try {
			const res = await apiFetch(`/api/autos/${id}`, { method: 'DELETE' });
			if (!res.ok) throw new Error('Error al eliminar');
			await cargarDatos();
		} catch (e) {
			alert('Error al eliminar el auto');
		}
	}

	const editingAuto = $derived(
		editingId ? autos.find((a: any) => a.id === editingId) ?? null : null
	);

	const tieneAutos = $derived(autos.length >= 4);

	async function nuevaCompetencia() {
		genLoading = true;
		error = '';
		success = '';
		try {
			const res = await apiFetch(`/api/categorias/${categoriaId}/competencias?rondas=${rondas}`, { method: 'POST' });
			if (!res.ok) {
				const text = await res.text();
				try { const d = JSON.parse(text); error = d.error || text; }
				catch { error = text || `Error ${res.status}`; }
				return;
			}
			// Recargar solo competencias
			const compRes = await fetch(`/api/categorias/${categoriaId}/competencias`);
			if (compRes.ok) {
				competencias = await compRes.json();
				success = '✅ Competencia creada con éxito';
				setTimeout(() => success = '', 3000);
			}
		} catch (e: any) { error = `Error: ${e.message || 'desconocido'}`; }
		finally { genLoading = false; }
	}

	function estadoLabel(estado: string): string {
		if (estado === 'abierta') return '🔴 En curso';
		if (estado === 'finalizada') return '🏁 Finalizada';
		return estado;
	}
</script>

<div class="header">
	<div>
		<a href="/admin/categorias" class="back-link">← Categorías</a>
		<h1>{categoria?.nombre ?? 'Cargando...'}</h1>
	</div>
	<button class="btn btn-primary" onclick={() => { editingId = null; showForm = true; }}>
		+ Nuevo Auto
	</button>
</div>

{#if error}
	<div class="error" in:fade>{error}</div>
{/if}
{#if success}
	<div class="success" in:fade>{success}</div>
{/if}

{#if showForm}
	<AutoForm
		auto={editingAuto}
		categoriaId={categoriaId}
		onclose={handleClose}
		onsaved={handleSaved}
	/>
{/if}

<div class="table-container">
	{#if loading}
		<p class="loading">Cargando...</p>
	{:else}
		<table>
			<thead>
				<tr>
					<th>#</th>
					<th>Nombre</th>
					<th>Creador</th>
					<th>Edad</th>
					<th>Foto</th>
					<th>Acciones</th>
				</tr>
			</thead>
			<tbody>
				{#each autos as auto (auto.id)}
					<tr>
						<td class="numero">{auto.numero}</td>
						<td>{auto.nombre}</td>
						<td>{auto.creador}</td>
						<td>{auto.edad}</td>
						<td>
							{#if auto.foto_url}
								<img src={auto.foto_url} alt={auto.nombre} class="auto-foto" />
							{:else}
								<span class="no-foto">—</span>
							{/if}
						</td>
						<td class="actions">
							<button class="btn btn-sm" onclick={() => handleEdit(auto.id)}>Editar</button>
							<button class="btn btn-sm btn-danger" onclick={() => handleDelete(auto.id)}>Eliminar</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="6" class="empty">No hay autos registrados en esta categoría</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- ═══ COMPETENCIAS ═══ -->
<section class="comp-section" in:fade>
	<div class="comp-header">
		<h2>🏁 Competencias</h2>
		{#if tieneAutos}
			<div class="comp-create">
				<label class="rondas-label">
					Rondas:
					<input type="number" bind:value={rondas} min="1" max="10" class="rondas-input" />
				</label>
				<button class="btn btn-primary" onclick={nuevaCompetencia} disabled={genLoading}>
					{genLoading ? 'Creando...' : '➕ Nueva Competencia'}
				</button>
			</div>
		{/if}
	</div>

	{#if !tieneAutos}
		<p class="comp-hint">Se necesitan al menos 4 autos para crear una competencia</p>
	{/if}

	{#if compLoading}
		<p class="comp-hint">Cargando competencias...</p>
	{:else if competencias.length === 0 && tieneAutos}
		<p class="comp-hint">Todavía no hay competencias. Creá la primera para empezar las carreras.</p>
	{/if}

	{#if competencias.length > 0}
		<div class="comp-list">
			{#each competencias as comp (comp.id)}
				<a href="/admin/categorias/{categoriaId}/fixture?competencia={comp.id}" class="comp-card" class:finalizada={comp.estado === 'finalizada'}>
					<div class="comp-info">
						<span class="comp-nombre">{comp.nombre}</span>
						<span class="comp-meta">{comp.rondas} rondas · {comp.estado === 'finalizada' ? 'Solo lectura' : 'Editable'}</span>
					</div>
					<span class="comp-estado">{estadoLabel(comp.estado)}</span>
				</a>
			{/each}
		</div>
	{/if}
</section>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 2rem;
	}

	h1 {
		color: #f59e0b;
		font-size: 1.75rem;
		margin: 0.5rem 0 0 0;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.back-link {
		color: #64748b;
		text-decoration: none;
		font-size: 0.875rem;
	}

	.back-link:hover { color: #f59e0b; }

	.error {
		color: #fca5a5;
		background: rgba(220,38,38,0.15);
		padding: 0.75rem 1rem;
		border-radius: 0.25rem;
		margin-bottom: 1rem;
		font-size: 0.9rem;
	}

	.success {
		color: #a3e635;
		background: rgba(101,163,13,0.15);
		padding: 0.75rem 1rem;
		border-radius: 0.25rem;
		margin-bottom: 1rem;
		font-size: 0.9rem;
		font-weight: 600;
	}

	.loading {
		text-align: center;
		color: #64748b;
		padding: 3rem;
	}

	.btn {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.25rem;
		cursor: pointer;
		font-weight: 600;
		transition: all 0.2s;
	}

	.btn-primary { background: #f59e0b; color: #0f172a; }
	.btn-primary:hover { background: #d97706; }
	.btn-sm { padding: 0.25rem 0.75rem; font-size: 0.875rem; background: #334155; color: #e2e8f0; }
	.btn-sm:hover { background: #475569; }
	.btn-danger { background: #dc2626; color: white; }
	.btn-danger:hover { background: #b91c1c; }

	.btn-fixture {
		display: inline-block;
		padding: 0.5rem 1rem;
		background: #dc2626;
		color: white;
		text-decoration: none;
		border-radius: 0.25rem;
		font-weight: 600;
	}

	.btn-fixture:hover { background: #b91c1c; }

	.table-container {
		background: #1e293b;
		border-radius: 0.5rem;
		overflow-x: auto;
		-webkit-overflow-scrolling: touch;
		border: 1px solid #334155;
	}

	@media (max-width: 640px) {
		.actions { flex-wrap: wrap; }
		.actions .btn { font-size: 0.75rem; padding: 0.2rem 0.5rem; }
		td, th { padding: 0.5rem; }
		.auto-foto { width: 32px; height: 32px; }
	}

	table { width: 100%; border-collapse: collapse; }

	th {
		background: #0f172a;
		color: #f59e0b;
		padding: 1rem;
		text-align: left;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	td { padding: 1rem; border-top: 1px solid #334155; }
	.numero { font-weight: bold; color: #f59e0b; font-size: 1.25rem; }
	.auto-foto { width: 48px; height: 48px; object-fit: cover; border-radius: 0.25rem; border: 1px solid #334155; }
	.no-foto { color: #475569; }
	.actions { display: flex; gap: 0.5rem; }
	.empty { text-align: center; color: #64748b; padding: 3rem; font-style: italic; }
	.nav-links { margin-top: 1.5rem; }

	/* Competencias */
	.comp-section { margin-top: 2rem; padding-top: 2rem; border-top: 2px solid var(--orange); }

	.comp-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; flex-wrap: wrap; margin-bottom: 1rem; }
	.comp-header h2 { color: var(--orange); font-size: 1.25rem; margin: 0; }

	.comp-create { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }

	.rondas-label { display: flex; align-items: center; gap: 0.3rem; color: var(--text-dim); font-size: 0.85rem; }
	.rondas-input { width: 60px; padding: 0.35rem; border: 1px solid var(--border-color); border-radius: 0.25rem; background: var(--arcade-black); color: var(--text-primary); text-align: center; }
	.rondas-input:focus { outline: none; border-color: var(--orange); }

	.comp-hint { color: var(--text-dim); font-size: 0.9rem; text-align: center; padding: 2rem; }

	.comp-list { display: flex; flex-direction: column; gap: 0.5rem; }

	.comp-card {
		display: flex; justify-content: space-between; align-items: center;
		background: var(--arcade-dark); border: 1px solid var(--border-color);
		border-radius: 0.5rem; padding: 1rem 1.25rem; text-decoration: none;
		transition: all 0.2s;
	}
	.comp-card:hover { border-color: var(--orange); transform: translateX(4px); }
	.comp-card.finalizada { opacity: 0.7; }
	.comp-card.finalizada:hover { border-color: var(--text-dim); opacity: 1; }

	.comp-info { display: flex; flex-direction: column; gap: 0.15rem; }
	.comp-nombre { color: var(--orange); font-weight: 600; }
	.comp-meta { color: var(--text-dim); font-size: 0.8rem; }
	.comp-estado { font-size: 0.85rem; font-weight: 600; white-space: nowrap; }
</style>

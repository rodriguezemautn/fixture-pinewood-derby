<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import CategoriaForm from '$lib/components/CategoriaForm.svelte';
	import { apiFetch } from '$lib/api';

	let categorias = $state<any[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showForm = $state(false);
	let editingId = $state<string | null>(null);

	onMount(() => {
		cargarCategorias();
	});

	async function cargarCategorias() {
		loading = true;
		error = '';
		try {
			const res = await fetch('/api/categorias');
			if (!res.ok) throw new Error('Error al cargar');
			categorias = await res.json();
		} catch (e) {
			error = 'No se pudieron cargar las categorías. ¿El backend está corriendo?';
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
		await cargarCategorias();
	}

	async function handleDelete(id: string) {
		if (!confirm('¿Eliminar esta categoría?')) return;
		try {
			const res = await apiFetch(`/api/categorias/${id}`, { method: 'DELETE' });
			if (!res.ok) throw new Error('Error al eliminar');
			await cargarCategorias();
		} catch (e) {
			alert('Error al eliminar la categoría');
		}
	}

	const editingCategoria = $derived(
		editingId ? categorias.find((c: any) => c.id === editingId) ?? null : null
	);
</script>

<div class="header">
	<h1>Categorías</h1>
	<button class="btn btn-primary" onclick={() => { editingId = null; showForm = true; }}>
		+ Nueva Categoría
	</button>
</div>

{#if error}
	<div class="error" in:fade>{error}</div>
{/if}

{#if showForm}
	<CategoriaForm
		categoria={editingCategoria}
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
					<th>Nombre</th>
					<th>Edad Mín</th>
					<th>Edad Máx</th>
					<th>Acciones</th>
				</tr>
			</thead>
			<tbody>
				{#each categorias as cat (cat.id)}
					<tr>
						<td>{cat.nombre}</td>
						<td>{cat.edad_min}</td>
						<td>{cat.edad_max}</td>
						<td class="actions">
							<a href="/admin/categorias/{cat.id}/autos" class="btn btn-sm btn-racing">🏎️ Autos</a>
							<a href="/admin/categorias/{cat.id}/fixture" class="btn btn-sm btn-fixture">🏁 Fixture</a>
							<button class="btn btn-sm" onclick={() => handleEdit(cat.id)}>Editar</button>
							<button class="btn btn-sm btn-danger" onclick={() => handleDelete(cat.id)}>Eliminar</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4" class="empty">No hay categorías registradas</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h1 {
		color: var(--orange);
		font-size: 2rem;
		margin: 0;
		text-transform: uppercase;
	}

	.error {
		color: #fca5a5;
		background: rgba(220,38,38,0.15);
		padding: 0.75rem 1rem;
		margin-bottom: 1rem;
		font-size: 0.9rem;
	}

	.loading {
		text-align: center;
		color: var(--text-dim);
		padding: 3rem;
	}

	.table-container {
		background: var(--arcade-surface);
		border: 1px solid var(--border-color);
		overflow: hidden;
	}

	table { width: 100%; border-collapse: collapse; }

	th {
		background: var(--arcade-dark);
		color: var(--orange);
		padding: 1rem;
		text-align: left;
		font-family: 'VT323', monospace;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
	}

	td {
		padding: 1rem;
		border-top: 1px solid var(--border-color);
		color: var(--text-body);
	}

	.actions { display: flex; gap: 0.5rem; }

	.empty {
		text-align: center;
		color: var(--text-dim);
		padding: 3rem;
		font-style: italic;
	}
</style>

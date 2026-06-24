<script lang="ts">
	import CategoriaForm from '$lib/components/CategoriaForm.svelte';
	import { apiFetch } from '$lib/api';

	let { data } = $props();
	let { categorias } = $derived(data);

	let showForm = $state(false);
	let editingId = $state<string | null>(null);

	function handleEdit(id: string) {
		editingId = id;
		showForm = true;
	}

	function handleClose() {
		showForm = false;
		editingId = null;
	}

	function handleSaved() {
		showForm = false;
		editingId = null;
		window.location.reload();
	}

	async function handleDelete(id: string) {
		if (!confirm('¿Eliminar esta categoría?')) return;
		await apiFetch(`/api/categorias/${id}`, { method: 'DELETE' });
		window.location.reload();
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

{#if showForm}
	<CategoriaForm
		categoria={editingCategoria}
		onclose={handleClose}
		onsaved={handleSaved}
	/>
{/if}

<div class="table-container">
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
</div>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h1 {
		color: #f59e0b;
		font-size: 2rem;
		margin: 0;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.btn {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.25rem;
		cursor: pointer;
		font-weight: 600;
		transition: all 0.2s;
	}

	.btn-primary {
		background: #f59e0b;
		color: #0f172a;
	}

	.btn-primary:hover {
		background: #d97706;
	}

	.btn-sm {
		padding: 0.25rem 0.75rem;
		font-size: 0.875rem;
		background: #334155;
		color: #e2e8f0;
	}

	.btn-sm:hover {
		background: #475569;
	}

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background: #b91c1c;
	}

	.btn-racing {
		background: #059669;
		color: white;
	}

	.btn-racing:hover {
		background: #047857;
	}

	.btn-fixture {
		background: #dc2626;
		color: white;
	}

	.btn-fixture:hover {
		background: #b91c1c;
	}

	.table-container {
		background: #1e293b;
		border-radius: 0.5rem;
		overflow: hidden;
		border: 1px solid #334155;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th {
		background: #0f172a;
		color: #f59e0b;
		padding: 1rem;
		text-align: left;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	td {
		padding: 1rem;
		border-top: 1px solid #334155;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
	}

	.empty {
		text-align: center;
		color: #64748b;
		padding: 3rem;
		font-style: italic;
	}
</style>

<script lang="ts">
	import AutoForm from '$lib/components/AutoForm.svelte';
	import { apiFetch } from '$lib/api';

	let { data } = $props();
	let { autos, categoria, categoriaId } = $derived(data);

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
		if (!confirm('¿Eliminar este auto?')) return;
		await apiFetch(`/api/autos/${id}`, { method: 'DELETE' });
		window.location.reload();
	}

	const editingAuto = $derived(
		editingId ? autos.find((a: any) => a.id === editingId) ?? null : null
	);
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

{#if showForm}
	<AutoForm
		auto={editingAuto}
		categoriaId={categoriaId}
		onclose={handleClose}
		onsaved={handleSaved}
	/>
{/if}

<div class="table-container">
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
</div>

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

	.back-link:hover {
		color: #f59e0b;
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

	.btn-primary:hover { background: #d97706; }

	.btn-sm {
		padding: 0.25rem 0.75rem;
		font-size: 0.875rem;
		background: #334155;
		color: #e2e8f0;
	}

	.btn-sm:hover { background: #475569; }

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover { background: #b91c1c; }

	.table-container {
		background: #1e293b;
		border-radius: 0.5rem;
		overflow: hidden;
		border: 1px solid #334155;
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

	td {
		padding: 1rem;
		border-top: 1px solid #334155;
	}

	.numero {
		font-weight: bold;
		color: #f59e0b;
		font-size: 1.25rem;
	}

	.auto-foto {
		width: 48px;
		height: 48px;
		object-fit: cover;
		border-radius: 0.25rem;
		border: 1px solid #334155;
	}

	.no-foto { color: #475569; }

	.actions { display: flex; gap: 0.5rem; }

	.empty {
		text-align: center;
		color: #64748b;
		padding: 3rem;
		font-style: italic;
	}
</style>

<script lang="ts">
	import { apiFetch } from '$lib/api';

	interface Categoria {
		id: string;
		nombre: string;
		edad_min: number;
		edad_max: number;
	}

	let { categoria = null, onclose, onsaved }: {
		categoria?: Categoria | null;
		onclose: () => void;
		onsaved: () => void;
	} = $props();

	let nombre = $state(categoria?.nombre ?? '');
	let edadMin = $state(categoria?.edad_min ?? 8);
	let edadMax = $state(categoria?.edad_max ?? 18);
	let error = $state('');
	let loading = $state(false);

	const isEditing = $derived(categoria !== null);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		const body = { nombre, edad_min: edadMin, edad_max: edadMax };
		const url = isEditing ? `/api/categorias/${categoria!.id}` : '/api/categorias';
		const method = isEditing ? 'PUT' : 'POST';

		try {
			const res = await apiFetch(url, {
				method,
				body: JSON.stringify(body)
			});

			if (!res.ok) {
				const data = await res.json();
				error = data.error || 'Error al guardar';
				return;
			}

			onsaved();
		} catch {
			error = 'Error de conexión';
		} finally {
			loading = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="modal-overlay" onclick={onclose} role="presentation">
	<!-- svelte-ignore a11y_interactive_supports_focus a11y_click_events_have_key_events -->
	<div class="modal" onclick={(e: MouseEvent) => e.stopPropagation()} role="dialog" aria-label="Formulario de categoría">
		<h2>{isEditing ? 'Editar' : 'Nueva'} Categoría</h2>

		<form onsubmit={handleSubmit}>
			<label>
				Nombre
				<input
					type="text"
					bind:value={nombre}
					required
					placeholder="Ej: Pre-Juveniles"
				/>
			</label>

			<div class="row">
				<label>
					Edad Mínima
					<input type="number" bind:value={edadMin} min="1" max="99" required />
				</label>
				<label>
					Edad Máxima
					<input type="number" bind:value={edadMax} min="1" max="99" required />
				</label>
			</div>

			{#if error}
				<p class="error">{error}</p>
			{/if}

			<div class="buttons">
				<button type="button" class="btn btn-cancel" onclick={onclose}>Cancelar</button>
				<button type="submit" class="btn btn-save" disabled={loading}>
					{loading ? 'Guardando...' : 'Guardar'}
				</button>
			</div>
		</form>
	</div>
</div>

<style>
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
	}

	.modal {
		background: var(--arcade-surface);
		border: 2px solid var(--orange-border);
		padding: 2rem;
		width: 100%;
		max-width: 500px;
		box-shadow: 4px 4px 0 0 var(--orange-border);
	}

	h2 {
		color: var(--orange);
		margin: 0 0 1.5rem 0;
		font-family: 'VT323', monospace;
		font-size: 1.3rem;
		letter-spacing: 0.05em;
	}

	form { display: flex; flex-direction: column; gap: 1rem; }

	label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		color: var(--text-dim);
		font-size: 0.875rem;
		flex: 1;
	}

	input {
		padding: 0.5rem;
		border: 1.5px solid var(--border-color);
		background: var(--arcade-black);
		color: var(--text-primary);
		font-family: 'VT323', monospace;
		font-size: 1.05rem;
		outline: none;
	}

	input:focus { border-color: var(--orange); box-shadow: 0 0 8px rgba(217, 119, 6, 0.15); }

	.row { display: flex; gap: 1rem; }

	.error { color: var(--red-race-light); font-size: 0.875rem; margin: 0; }

	.buttons { display: flex; gap: 0.5rem; justify-content: flex-end; margin-top: 1rem; }
</style>

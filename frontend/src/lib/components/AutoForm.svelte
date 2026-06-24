<script lang="ts">
	interface Auto {
		id: string;
		numero: number;
		nombre: string;
		creador: string;
		edad: number;
		foto_url: string;
	}

	let { auto = null, categoriaId, onclose, onsaved }: {
		auto?: Auto | null;
		categoriaId: string;
		onclose: () => void;
		onsaved: () => void;
	} = $props();

	let numero = $state(auto?.numero ?? 1);
	let nombre = $state(auto?.nombre ?? '');
	let creador = $state(auto?.creador ?? '');
	let edad = $state(auto?.edad ?? 8);
	let fotoUrl = $state(auto?.foto_url ?? '');
	let error = $state('');
	let loading = $state(false);

	const isEditing = $derived(auto !== null);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		const body = { numero, nombre, creador, edad, foto_url: fotoUrl };

		try {
			const res = isEditing
				? await fetch(`/api/autos/${auto!.id}`, {
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(body)
				})
				: await fetch(`/api/categorias/${categoriaId}/autos`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
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
	<div class="modal" onclick={(e: MouseEvent) => e.stopPropagation()} role="dialog" aria-label="Formulario de auto">
		<h2>{isEditing ? 'Editar' : 'Nuevo'} Auto</h2>

		<form onsubmit={handleSubmit}>
			<label>
				Número
				<input type="number" bind:value={numero} min="1" max="999" required />
			</label>

			<label>
				Nombre del Auto
				<input type="text" bind:value={nombre} required placeholder="Ej: Turbo Relámpago" />
			</label>

			<label>
				Creador
				<input type="text" bind:value={creador} required placeholder="Nombre del participante" />
			</label>

			<label>
				Edad
				<input type="number" bind:value={edad} min="1" max="99" required />
			</label>

			<label>
				Foto (URL)
				<input type="url" bind:value={fotoUrl} placeholder="https://ejemplo.com/foto.jpg" />
			</label>

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
		background: #1e293b;
		border: 1px solid #f59e0b;
		border-radius: 0.5rem;
		padding: 2rem;
		width: 100%;
		max-width: 500px;
	}

	h2 {
		color: #f59e0b;
		margin: 0 0 1.5rem 0;
	}

	form { display: flex; flex-direction: column; gap: 1rem; }

	label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		color: #94a3b8;
		font-size: 0.875rem;
	}

	input {
		padding: 0.5rem;
		border: 1px solid #334155;
		border-radius: 0.25rem;
		background: #0f172a;
		color: #e2e8f0;
		font-size: 1rem;
	}

	input:focus {
		outline: none;
		border-color: #f59e0b;
	}

	.error {
		color: #ef4444;
		font-size: 0.875rem;
		margin: 0;
	}

	.buttons {
		display: flex;
		gap: 0.5rem;
		justify-content: flex-end;
		margin-top: 1rem;
	}

	.btn {
		padding: 0.5rem 1.5rem;
		border: none;
		border-radius: 0.25rem;
		cursor: pointer;
		font-weight: 600;
		font-size: 0.875rem;
	}

	.btn-cancel { background: #334155; color: #e2e8f0; }
	.btn-save { background: #f59e0b; color: #0f172a; }
	.btn-save:disabled { opacity: 0.5; cursor: not-allowed; }
</style>

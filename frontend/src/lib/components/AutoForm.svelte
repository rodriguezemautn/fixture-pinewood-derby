<script lang="ts">
	import { apiFetch } from '$lib/api';

	interface Auto {
		id: string;
		numero: number;
		nombre: string;
		creador: string;
		edad: number;
		peso: number;
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
	let peso = $state(auto?.peso ?? 0);
	let fotoUrl = $state(auto?.foto_url ?? '');
	let fotoFile = $state<File | null>(null);
	let fotoPreview = $state<string | null>(null);
	let error = $state('');
	let loading = $state(false);

	const isEditing = $derived(auto !== null);

	function handleFileChange(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			fotoFile = input.files[0];
			fotoPreview = URL.createObjectURL(input.files[0]);
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		const body = { numero, nombre, creador, edad, peso, foto_url: fotoUrl };

		try {
			// Primero crear/actualizar el auto
			const res = isEditing
				? await apiFetch(`/api/autos/${auto!.id}`, {
					method: 'PUT',
					body: JSON.stringify(body)
				})
				: await apiFetch(`/api/categorias/${categoriaId}/autos`, {
					method: 'POST',
					body: JSON.stringify(body)
				});

			if (!res.ok) {
				const data = await res.json();
				error = data.error || 'Error al guardar';
				return;
			}

			const autoCreado = await res.json();

			// Si hay foto para subir, subirla
			if (fotoFile) {
				const formData = new FormData();
				formData.append('foto', fotoFile);

				const fotoRes = await apiFetch(`/api/autos/${autoCreado.id}/foto`, {
					method: 'POST',
					body: formData,
					headers: {} // dejar que fetch setee Content-Type con boundary
				});

				if (!fotoRes.ok) {
					error = 'Auto guardado pero la foto no pudo subirse';
					return;
				}
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
				Peso (gramos)
				<input type="number" bind:value={peso} min="0" max="5000" placeholder="Ej: 150" />
				<span class="file-hint">Peso del auto en gramos</span>
			</label>

			<label>
				Foto
				<input type="file" accept="image/jpeg,image/png,image/gif,image/webp" capture="environment" onchange={handleFileChange} />
				<span class="file-hint">Foto desde el celular o archivo</span>
			</label>

			{#if fotoPreview}
				<div class="preview">
					<img src={fotoPreview} alt="Preview" class="preview-img" />
				</div>
			{/if}

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

	h2 { color: #f59e0b; margin: 0 0 1.5rem 0; }
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

	input:focus { outline: none; border-color: #f59e0b; }
	input[type="file"] { font-size: 0.85rem; padding: 0.4rem; }

	.file-hint { color: #64748b; font-size: 0.75rem; }

	.preview { text-align: center; }
	.preview-img { max-width: 150px; max-height: 150px; border-radius: 0.25rem; border: 1px solid #334155; }

	.error { color: #ef4444; font-size: 0.875rem; margin: 0; }
	.buttons { display: flex; gap: 0.5rem; justify-content: flex-end; margin-top: 1rem; }

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

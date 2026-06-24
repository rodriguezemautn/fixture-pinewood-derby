<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';

	let pin = $state('');
	let error = $state('');
	let loading = $state(false);

	// Si ya está autenticado, redirigir
	if (browser && localStorage.getItem('auth_token')) {
		goto('/admin/categorias');
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ pin })
			});

			if (!res.ok) {
				const data = await res.json();
				error = data.error || 'PIN incorrecto';
				return;
			}

			const data = await res.json();
			localStorage.setItem('auth_token', data.token);
			localStorage.setItem('auth_role', data.role);
			goto('/admin/categorias');
		} catch {
			error = 'Error de conexión';
		} finally {
			loading = false;
		}
	}
</script>

<div class="login-page">
	<div class="login-card" in:fly={{ y: 20, duration: 400 }}>
		<div class="login-header">
			<img src="/assets/emblema.jpg" alt="Emblema" class="login-logo" />
			<h1>Fixture D15</h1>
			<p class="login-sub">Acceso Administrador</p>
		</div>

		<form onsubmit={handleSubmit}>
			<label>
				PIN de Administrador
				<input
					type="password"
					bind:value={pin}
					placeholder="Ingresá tu PIN"
					maxlength="10"
					autofocus
					required
				/>
			</label>

			{#if error}
				<p class="error" in:fade>{error}</p>
			{/if}

			<button type="submit" class="btn-login" disabled={loading || !pin}>
				{loading ? 'Ingresando...' : '🔐 Ingresar'}
			</button>
		</form>

		<a href="/" class="back-link">← Volver a la página principal</a>
	</div>
</div>

<style>
	.login-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: radial-gradient(ellipse at center, #1a1f2e 0%, #0a0f1a 70%);
		padding: 1rem;
	}

	.login-card {
		background: var(--racing-gray);
		border: 1px solid var(--racing-amber);
		border-radius: 0.75rem;
		padding: 2.5rem;
		width: 100%;
		max-width: 380px;
		text-align: center;
		box-shadow: 0 0 30px rgba(245,158,11,0.15);
	}

	.login-header { margin-bottom: 2rem; }

	.login-logo {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		border: 3px solid var(--racing-amber);
		margin-bottom: 1rem;
	}

	h1 {
		color: var(--racing-amber);
		font-size: 1.5rem;
		margin: 0;
	}

	.login-sub {
		color: var(--racing-text-dim);
		font-size: 0.9rem;
		margin: 0.5rem 0 0;
	}

	form { display: flex; flex-direction: column; gap: 1rem; }

	label {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		color: var(--racing-text-dim);
		font-size: 0.85rem;
		text-align: left;
	}

	input {
		padding: 0.75rem;
		border: 1px solid var(--racing-border);
		border-radius: 0.35rem;
		background: var(--racing-black);
		color: var(--racing-text);
		font-size: 1.25rem;
		letter-spacing: 0.5em;
		text-align: center;
	}

	input:focus {
		outline: none;
		border-color: var(--racing-amber);
		box-shadow: 0 0 10px rgba(245,158,11,0.2);
	}

	.error {
		color: #fca5a5;
		font-size: 0.85rem;
		margin: 0;
		background: rgba(220,38,38,0.15);
		padding: 0.5rem;
		border-radius: 0.25rem;
	}

	.btn-login {
		padding: 0.75rem;
		background: linear-gradient(135deg, var(--racing-amber), var(--racing-amber-light));
		color: var(--racing-black);
		font-weight: 800;
		font-size: 1rem;
		border: none;
		border-radius: 0.35rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-login:hover { transform: translateY(-2px); box-shadow: 0 4px 15px rgba(245,158,11,0.3); }
	.btn-login:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }

	.back-link {
		display: block;
		margin-top: 1.5rem;
		color: var(--racing-text-dim);
		text-decoration: none;
		font-size: 0.85rem;
	}

	.back-link:hover { color: var(--racing-amber); }
</style>

<script lang="ts">
	import { fade } from 'svelte/transition';
	import { browser } from '$app/environment';
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';

	let { children } = $props();

	// Auth guard: si no hay token, redirigir al login
	if (browser) {
		const token = localStorage.getItem('auth_token');
		if (!token) {
			goto('/login');
		}
	}

	function logout() {
		localStorage.removeItem('auth_token');
		localStorage.removeItem('auth_role');
		goto('/login');
	}

	let role = $state('');
	if (browser) {
		role = localStorage.getItem('auth_role') ?? '';
	}
</script>

<div class="admin-layout">
	<nav class="admin-nav">
		<div class="admin-nav-inner">
			<a href="/" class="admin-brand">
				<img src="/assets/derby_logo.jpg" alt="Derby" class="admin-logo" />
				<span>Fixture D15</span>
			</a>
			<div class="admin-links">
				<a href="/admin/categorias" class="admin-link">Categorías</a>
				<a href="/" class="admin-link">Ver Carrera</a>
			</div>
			<button class="btn-logout" onclick={logout}>Salir</button>
		</div>
		<div class="admin-nav-stripe"></div>
	</nav>

	<div class="admin-body" in:fade={{ duration: 200 }}>
		{@render children()}
	</div>
</div>

<style>
	.admin-layout {
		min-height: 100vh;
		background: var(--racing-black);
	}

	.admin-nav {
		background: linear-gradient(180deg, #0a0f1a, #0f172a);
		border-bottom: 2px solid var(--racing-amber);
	}

	.admin-nav-inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.6rem 1rem;
		max-width: 1200px;
		margin: 0 auto;
		gap: 1rem;
	}

	.admin-brand {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		text-decoration: none;
		color: var(--racing-amber);
		font-family: 'Black Ops One', sans-serif;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.admin-logo {
		height: 28px;
		width: auto;
		border-radius: 0.15rem;
	}

	.admin-links {
		display: flex;
		gap: 0.25rem;
	}

	.admin-link {
		padding: 0.4rem 0.75rem;
		color: var(--racing-text-dim);
		text-decoration: none;
		font-size: 0.85rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	.admin-link:hover {
		color: var(--racing-amber);
		background: rgba(245,158,11,0.1);
	}

	.btn-logout {
		padding: 0.35rem 0.75rem;
		background: transparent;
		border: 1px solid var(--racing-red);
		color: var(--racing-red);
		border-radius: 0.25rem;
		cursor: pointer;
		font-size: 0.8rem;
		font-weight: 600;
		transition: all 0.2s;
	}

	.btn-logout:hover {
		background: var(--racing-red);
		color: white;
	}

	.admin-nav-stripe {
		height: 2px;
		background: repeating-linear-gradient(
			90deg,
			var(--racing-amber) 0,
			var(--racing-amber) 6px,
			transparent 6px,
			transparent 12px
		);
	}

	.admin-body {
		padding: 1.5rem 1rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	@media (max-width: 640px) {
		.admin-brand span { display: none; }
		.admin-body { padding: 1rem 0.5rem; }
	}
</style>

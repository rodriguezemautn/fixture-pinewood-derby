<script lang="ts">
	import { fade } from 'svelte/transition';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
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

	function toggleTheme() {
		if (!browser) return;
		const current = document.documentElement.getAttribute('data-theme') || 'dark';
		const next = current === 'dark' ? 'light' : 'dark';
		document.documentElement.setAttribute('data-theme', next);
		localStorage.setItem('arcade-theme', next);
	}

	// Inicializar tema si ya estaba guardado
	if (browser) {
		const saved = localStorage.getItem('arcade-theme');
		if (saved === 'light' || saved === 'dark') {
			document.documentElement.setAttribute('data-theme', saved);
		}
	}

	// Breadcrumbs
	const breadcrumbs = $derived.by(() => {
		const path = $page.url.pathname;
		const parts = path.split('/').filter(Boolean);
		const crumbs: {label: string; href: string}[] = [];
		let accumulated = '';

		for (const part of parts) {
			accumulated += '/' + part;
			let label = part;
			if (part === 'admin') label = 'Admin';
			else if (part === 'categorias') label = 'Categorías';
			else if (part === 'autos') label = 'Autos';
			else if (part === 'fixture') label = 'Fixture';
			// IDs de UUID los mostramos truncados
			else if (/^[0-9a-f-]{36}$/.test(part)) label = 'Detalle';
			crumbs.push({ label, href: accumulated });
		}
		return crumbs;
	});
</script>

<div class="admin-layout">
	<nav class="admin-nav">
		<div class="admin-nav-inner">
			<a href="/" class="admin-brand">
				<img src="/assets/derby_logo.jpg" alt="Derby" class="admin-logo" />
				<span class="font-pixel" style="font-size: 0.7rem;">Fixture D15</span>
			</a>
			<div class="admin-links">
				<a href="/admin/categorias" class="admin-link">Categorías</a>
				<a href="/carreras" class="admin-link">Ver Carreras</a>
			</div>
			<button class="btn-theme" onclick={toggleTheme} title="Cambiar tema">🎨</button>
			<button class="btn-logout" onclick={logout}>Salir</button>
		</div>
		<div class="admin-nav-stripe"></div>
	</nav>

	<div class="admin-breadcrumbs">
		{#each breadcrumbs as crumb, i}
			{#if i > 0}<span class="crumb-sep">›</span>{/if}
			<a href={crumb.href} class="crumb">{crumb.label}</a>
		{/each}
	</div>

	<div class="admin-body" in:fade={{ duration: 200 }}>
		{@render children()}
	</div>
</div>

<style>
	.admin-layout {
		min-height: 100vh;
		background: var(--arcade-black);
	}

	.admin-nav {
		background: linear-gradient(180deg, var(--arcade-black), var(--arcade-dark));
		border-bottom: 2px solid var(--orange);
		position: relative;
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
		color: var(--orange);
		text-transform: uppercase;
	}

	.admin-logo {
		height: 28px;
		width: auto;
		image-rendering: pixelated;
		border: 1px solid var(--orange-border);
	}

	.admin-links {
		display: flex;
		gap: 0.25rem;
	}

	.admin-link {
		padding: 0.4rem 0.75rem;
		color: var(--text-dim);
		text-decoration: none;
		font-family: 'VT323', monospace;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		border: 1.5px solid transparent;
		transition: all 0.15s ease;
	}

	.admin-link:hover {
		color: var(--orange);
		border-color: var(--orange-border);
		box-shadow: 1.5px 1.5px 0 0 var(--orange-border);
	}

	.btn-theme {
		padding: 0.35rem 0.5rem;
		background: transparent;
		border: 1.5px solid var(--orange-border);
		color: var(--orange);
		cursor: pointer;
		font-size: 0.85rem;
		line-height: 1;
		box-shadow: 1.5px 1.5px 0 0 var(--orange-border);
		transition: all 0.15s ease;
	}
	.btn-theme:hover {
		transform: translate(-1px, -1px);
		box-shadow: 2px 2px 0 0 var(--orange);
		border-color: var(--orange);
	}

	.btn-logout {
		padding: 0.35rem 0.7rem;
		background: transparent;
		border: 1.5px solid var(--red-race);
		color: var(--red-race);
		cursor: pointer;
		font-family: 'VT323', monospace;
		font-size: 0.95rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		box-shadow: 1.5px 1.5px 0 0 var(--red-race);
		transition: all 0.15s ease;
	}

	.btn-logout:hover {
		background: var(--red-race);
		color: white;
		transform: translate(-1px, -1px);
		box-shadow: 2px 2px 0 0 var(--red-race);
	}

	.admin-nav-stripe {
		height: 2px;
		background: repeating-linear-gradient(
			90deg,
			var(--orange) 0,
			var(--orange) 6px,
			transparent 6px,
			transparent 12px
		);
	}

	.admin-breadcrumbs {
		display: flex; align-items: center; gap: 0.3rem;
		max-width: 1200px; margin: 0.5rem auto 0; padding: 0 1rem;
		font-size: 0.85rem;
		font-family: 'VT323', monospace;
		color: var(--text-dim);
	}

	.crumb { color: var(--text-dim); text-decoration: none; transition: color 0.15s; }
	.crumb:hover { color: var(--orange); }
	.crumb-sep { color: var(--border-color); font-size: 1rem; }

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

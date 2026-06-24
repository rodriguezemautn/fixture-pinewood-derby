<script lang="ts">
	let { children } = $props();
</script>

<div class="admin-layout">
	<nav class="admin-nav">
		<a href="/admin/categorias" class="nav-link">🏎️ Categorías</a>
		<a href="/" class="nav-link">← Ver Carrera</a>
	</nav>
	<main class="admin-content">
		{@render children()}
	</main>
</div>

<style>
	.admin-layout {
		min-height: 100vh;
		background: #0f172a;
		color: #e2e8f0;
	}

	.admin-nav {
		display: flex;
		gap: 1rem;
		padding: 1rem 2rem;
		background: #1e293b;
		border-bottom: 2px solid #f59e0b;
	}

	.nav-link {
		color: #f59e0b;
		text-decoration: none;
		font-weight: 600;
		padding: 0.5rem 1rem;
		border-radius: 0.25rem;
		transition: background 0.2s;
	}

	.nav-link:hover {
		background: #334155;
	}

	.admin-content {
		padding: 2rem;
		max-width: 1200px;
		margin: 0 auto;
	}
</style>

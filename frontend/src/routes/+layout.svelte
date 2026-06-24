<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Fixture — Pinewood Derby D15</title>
</svelte:head>

<!-- ════ HEADER RACING ════ -->
<header class="racing-header">
	<div class="header-inner">
		<a href="/" class="logo-link">
			<img
				src="/assets/derby_logo.jpg"
				alt="Pinewood Derby"
				class="logo-img"
			/>
			<span class="logo-text">Fixture D15</span>
		</a>
		<nav class="header-nav">
			<a href="/admin/categorias" class="nav-btn">Admin</a>
		</nav>
	</div>
	<div class="header-stripe"></div>
</header>

<!-- ════ MAIN CONTENT ════ -->
<main class="main-content">
	{@render children()}
</main>

<!-- ════ FOOTER ════ -->
<footer class="racing-footer">
	<div class="footer-stripe"></div>
	<p>Destacamento 15 — Iglesia Betel &bull; Pinewood Derby</p>
</footer>

<style>
	/* ─── Header ─────────────────────────────── */
	.racing-header {
		background: linear-gradient(180deg, #0a0f1a 0%, #0f172a 100%);
		border-bottom: 2px solid var(--racing-amber);
		position: sticky;
		top: 0;
		z-index: 50;
	}

	.header-inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.logo-link {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		text-decoration: none;
	}

	.logo-img {
		height: 40px;
		width: auto;
		border-radius: 0.25rem;
	}

	.logo-text {
		font-family: 'Black Ops One', sans-serif;
		font-size: 1.25rem;
		color: var(--racing-amber);
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.header-nav {
		display: flex;
		gap: 0.5rem;
	}

	.nav-btn {
		padding: 0.4rem 1rem;
		border: 1px solid var(--racing-amber);
		border-radius: 0.25rem;
		color: var(--racing-amber);
		text-decoration: none;
		font-weight: 600;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		transition: all 0.2s;
	}

	.nav-btn:hover {
		background: var(--racing-amber);
		color: var(--racing-black);
	}

	.header-stripe {
		height: 3px;
		background: repeating-linear-gradient(
			90deg,
			var(--racing-amber) 0,
			var(--racing-amber) 8px,
			transparent 8px,
			transparent 16px
		);
	}

	/* ─── Main ───────────────────────────────── */
	.main-content {
		min-height: calc(100vh - 120px);
	}

	/* ─── Footer ─────────────────────────────── */
	.racing-footer {
		background: var(--racing-dark);
		text-align: center;
		padding: 1.5rem 1rem;
		color: var(--racing-text-dim);
		font-size: 0.8rem;
	}

	.footer-stripe {
		height: 2px;
		background: repeating-linear-gradient(
			90deg,
			var(--racing-red) 0,
			var(--racing-red) 6px,
			transparent 6px,
			transparent 12px
		);
		margin-bottom: 1rem;
	}

	@media (max-width: 640px) {
		.logo-text { font-size: 1rem; }
		.logo-img { height: 32px; }
	}
</style>

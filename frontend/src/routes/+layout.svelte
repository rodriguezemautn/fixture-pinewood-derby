<script lang="ts">
	import { onMount } from 'svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	let theme = $state('dark');

	onMount(() => {
		const saved = localStorage.getItem('arcade-theme');
		if (saved === 'light' || saved === 'dark') {
			theme = saved;
		}
		applyTheme(theme);
	});

	function toggleTheme() {
		theme = theme === 'dark' ? 'light' : 'dark';
		localStorage.setItem('arcade-theme', theme);
		applyTheme(theme);
	}

	function applyTheme(t: string) {
		document.documentElement.setAttribute('data-theme', t);
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Fixture — Pinewood Derby D15</title>
</svelte:head>

<!-- ════ HEADER ARCADE RACING ════ -->
<header class="arcade-header">
	<div class="header-inner">
		<a href="/" class="logo-link">
			<img
				src="/assets/derby_logo.jpg"
				alt="Pinewood Derby"
				class="logo-img"
			/>
			<span class="logo-text font-pixel">Fixture D15</span>
		</a>
		<nav class="header-nav">
			<a href="/admin/categorias" class="nav-btn">Admin</a>
			<button class="theme-btn" onclick={toggleTheme} title="Cambiar tema">
				{theme === 'dark' ? '☀️' : '🌙'}
			</button>
		</nav>
	</div>
	<div class="header-stripe"></div>
	<div class="header-scanline"></div>
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
	/* ─── Header Arcade ───────────────────────── */
	.arcade-header {
		background: linear-gradient(180deg, var(--arcade-black) 0%, var(--arcade-dark) 100%);
		border-bottom: 2px solid var(--orange);
		position: sticky;
		top: 0;
		z-index: 50;
		position: relative;
		overflow: hidden;
	}

	.header-scanline {
		position: absolute;
		inset: 0;
		pointer-events: none;
		background: repeating-linear-gradient(
			0deg,
			transparent 0,
			transparent 2px,
			rgba(0, 0, 0, 0.12) 2px,
			rgba(0, 0, 0, 0.12) 4px
		);
	}

	.header-inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		max-width: 1200px;
		margin: 0 auto;
		position: relative;
		z-index: 1;
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
		image-rendering: pixelated;
		border: 1.5px solid var(--orange-border);
	}

	.logo-text {
		font-size: 0.85rem;
		color: var(--orange);
		text-transform: uppercase;
	}

	.header-nav {
		display: flex;
		gap: 0.5rem;
	}

	.nav-btn {
		padding: 0.4rem 0.9rem;
		border: 2px solid var(--orange-border);
		color: var(--orange);
		text-decoration: none;
		font-family: 'VT323', monospace;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		transition: all 0.15s ease;
		box-shadow: 2px 2px 0 0 var(--orange-border);
	}

	.nav-btn:hover {
		background: var(--orange);
		color: var(--arcade-black);
		box-shadow: 2px 2px 0 0 var(--orange);
		transform: translate(-1px, -1px);
	}

	.theme-btn {
		padding: 0.35rem 0.6rem;
		border: 2px solid var(--orange-border);
		background: transparent;
		cursor: pointer;
		font-size: 1.1rem;
		line-height: 1;
		box-shadow: 2px 2px 0 0 var(--orange-border);
		transition: all 0.15s ease;
	}

	.theme-btn:hover {
		transform: translate(-1px, -1px);
		box-shadow: 2px 2px 0 0 var(--orange);
		border-color: var(--orange);
	}

	.header-stripe {
		height: 3px;
		background: repeating-linear-gradient(
			90deg,
			var(--orange) 0,
			var(--orange) 8px,
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
		background: var(--arcade-dark);
		text-align: center;
		padding: 1.5rem 1rem;
		color: var(--text-dim);
		font-size: 0.8rem;
		border-top: 1px solid var(--border-color);
	}

	.footer-stripe {
		height: 2px;
		background: repeating-linear-gradient(
			90deg,
			var(--red-race) 0,
			var(--red-race) 6px,
			transparent 6px,
			transparent 12px
		);
		margin-bottom: 1rem;
	}

	@media (max-width: 640px) {
		.logo-text { font-size: 0.7rem; }
		.logo-img { height: 32px; }
	}
</style>

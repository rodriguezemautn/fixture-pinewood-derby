# Proposal: Estilo Racing

## Intent
Transformar el frontend con temática gamer/de carreras usando los assets del proyecto (logos D15, emblemas, autos).

## Scope
**In**: Landing page con identidad racing, admin con estilo consistente, animaciones, transiciones Svelte 5, layout general con header/nav racing, favicon/fonts, diseño mobile-first responsivo.
**Out**: Lógica de negocio, funcionalidad nueva, backend.

## Approach
Sobreescribir +layout.svelte global con header racing (logo + navegación), +page.svelte como dashboard con hero, animaciones con transiciones nativas de Svelte 5, paleta racing (dark #0f172a, amber #f59e0b, red accent), tipografía heading impactante.

## Files
- `frontend/src/routes/+layout.svelte` (modify)
- `frontend/src/routes/+page.svelte` (modify)
- `frontend/src/routes/+layout.css` (modify)
- `frontend/src/routes/admin/+layout.svelte` (modify)
- `frontend/src/app.html` (modify — lang, meta, tema)

## Success Criteria
- Landing page con hero racing usando logo y emblema
- Navegación global con estilo consistente
- Transiciones suaves entre páginas
- Diseño mobile-first responsivo
- Assets locales funcionando en build

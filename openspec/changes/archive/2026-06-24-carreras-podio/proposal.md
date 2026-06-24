# Proposal: Carreras y Podio

## Intent
Implementar visualización gráfica del fixture, podios animados por carrera, y ventana especial para la gran final (RF6, RF7, RF8).

## Scope
**In**: Vista gráfica del fixture con estado de heats, podio animado al completar heat, ventana especial para podio final con celebración, vista pública del fixture para proyección.
**Out**: Autenticación, scheduling (ya implementado).

## Capabilities
- New: `visualizacion-carreras` — Podios animados + vista fixture gráfica

## Approach
Frontend pesado: Svelte 5 transiciones + CSS animations. Ruta pública `/carreras/{categoriaId}` para proyección (baja resolución, mobile-first). Podium reveal con animaciones escalonadas. Final con confetti/fuegos artificiales CSS.

## Files
- `frontend/src/routes/carreras/[id]/+page.svelte` — Vista pública del fixture
- `frontend/src/lib/components/Podium.svelte` — Componente podio animado
- `frontend/src/lib/components/Celebration.svelte` — Animación celebración final
- Modificaciones en admin fixture page para mostrar podio

## Success Criteria
- [ ] Vista pública del fixture con heats y posiciones
- [ ] Podio animado al hacer clic en resultado de heat
- [ ] Ventana de celebración para carrera final
- [ ] Diseño mobile-first, apto para proyección

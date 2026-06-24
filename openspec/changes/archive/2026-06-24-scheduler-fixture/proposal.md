# Proposal: Scheduler de Fixture

## Intent
Generar fixture de carreras óptimo y equitativo para determinar los 4 mejores autos de cada categoría para la carrera final (RF3, RF4, RF5).

## Scope
**In**: Algoritmo de scheduler, API para gestionar fixture, registro de resultados, cálculo de posiciones, selección top 4 para final, frontend de fixture y posiciones.
**Out**: Animaciones de carrera, podio emocionante (RF7-RF8 — próximo cambio), foto finish.

## Capabilities
- New: `scheduler-fixture` — Algoritmo + API + Frontend fixture

## Approach
Round-robin por grupos con R rondas de clasificación. Heats de 4 autos. Puntos por posición (1°=5, 2°=3, 3°=2, 4°=1, DNS=0). Top 4 acumulado → final. Service puro en Go sin dependencias.

## Affected
| Area | Impact |
|------|--------|
| `backend/internal/domain/` | New types: Fixture, Resultado, Standing |
| `backend/internal/service/fixture.go` | New: scheduler algorithm |
| `backend/internal/repository/` | FixtureRepository interface + SQLite |
| `backend/internal/handler/fixture.go` | New: fixture API endpoints |
| `backend/internal/database/sqlite.go` | New migrations: fixtures, resultados |

## Success Criteria
- [ ] Genera fixture para N autos con R rondas
- [ ] Cada auto corre R veces antes de la final
- [ ] Puntos calculados correctamente
- [ ] Top 4 seleccionados para la final
- [ ] API REST: crear fixture, registrar resultado, ver posiciones
- [ ] Tests > 70% coverage en el algoritmo

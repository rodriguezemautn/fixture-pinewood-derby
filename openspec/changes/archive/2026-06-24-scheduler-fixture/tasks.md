# Tasks: Scheduler de Fixture

## Phase 1: Domain + Algorithm
- [ ] 1.1 Domain types: Fixture, Heat, Resultado, Standing
- [ ] 1.2 **[TDD]** Scheduler algorithm: generar heats, calcular puntos, seleccionar top 4
- [ ] 1.3 Tests para algoritmo (varios N, R, empates)

## Phase 2: Persistence
- [ ] 2.1 Migraciones DB: tabla fixtures, heats, resultados
- [ ] 2.2 FixtureRepository interface + SQLite impl
- [ ] 2.3 Tests repository

## Phase 3: Service + Handler
- [ ] 3.1 FixtureService orquestador
- [ ] 3.2 Handler REST + tests
- [ ] 3.3 main.go wiring

## Phase 4: Frontend
- [ ] 4.1 Página fixture con heats y posiciones
- [ ] 4.2 Registro de resultados (orden de llegada)
- [ ] 4.3 Tabla de posiciones en vivo
- [ ] 4.4 Botón generar final

## Phase 5: Verify
- [ ] 5.1 go test ./... pasa
- [ ] 5.2 go vet
- [ ] 5.3 make build

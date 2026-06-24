# Tasks: Registro de Autos

## Phase 1: Infraestructura

- [x] 1.1 Migración DB: tabla autos con FK a categorías
- [x] 1.2 Agregar timestamps (CreatedAt, UpdatedAt) a domain.Auto
- [x] 1.3 Extender AutoRepository interface (Update, ExistsByNumero)

## Phase 2: Backend Core (TDD)

- [x] 2.1-2.2 Repository tests + impl (5 tests)
- [x] 2.3-2.4 Service tests + impl (6 tests)
- [x] 2.5-2.6 Handler tests + impl (7 tests)

## Phase 3: Wiring + Frontend

- [x] 3.1 Actualizar `main.go` con auto handler
- [x] 3.2 Ruta frontend `/admin/categorias/[id]/autos` con listado
- [x] 3.3 Componente `AutoForm.svelte` (crear/editar con foto)
- [x] 3.4 Navegación desde listado de categorías

## Phase 4: Testing & Verify

- [x] 4.1 `go test ./...` → pasa (handler 69.7%, service 82.5%, repo 80.5%)
- [x] 4.2 `go vet ./...` → pasa
- [x] 4.3 `make build` → binario + frontend compilados

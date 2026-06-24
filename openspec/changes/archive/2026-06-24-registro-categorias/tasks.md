# Tasks: Registro de Categorías

## Phase 1: Infraestructura DB

- [x] 1.1 Agregar dependencias: `modernc.org/sqlite` + `google/uuid`
- [x] 1.2 Agregar timestamps (CreatedAt, UpdatedAt) a domain.Categoria
- [x] 1.3 Crear `internal/database/sqlite.go`: conexión + migración tabla categorías
- [x] 1.4 Actualizar `cmd/api/main.go`: inicializar DB y pasar dependencias

## Phase 2: Backend Core (TDD)

- [x] 2.1 Extender `CategoriaRepository` interface con FindByID, Update, Delete
- [x] 2.2-2.3 Repository SQLite concreto + tests (4 tests, 80.5% coverage)
- [x] 2.4-2.5 Service con validaciones + tests (7 tests, 86.7% coverage)
- [x] 2.6 Extender interface CategoriaService con Update, Delete
- [x] 2.7-2.8 Handler REST completo + tests (7 tests, 72.4% coverage)

## Phase 3: Frontend

- [x] 3.1 Layout admin con navegación
- [x] 3.2 Ruta `/admin/categorias` con listado SSR
- [x] 3.3 Componente `CategoriaForm.svelte` (crear/editar modal)
- [x] 3.4 Conexión fetch con API REST

## Phase 4: Testing & Verify

- [x] 4.1 `go test ./...` → pasa (handler 72.4%, service 86.7%, repo 80.5%)
- [x] 4.2 `go vet ./...` → pasa
- [x] 4.3 `make build` → binario + frontend compilados
- [x] 4.4 CRUD manual contra API (verificado con tests)

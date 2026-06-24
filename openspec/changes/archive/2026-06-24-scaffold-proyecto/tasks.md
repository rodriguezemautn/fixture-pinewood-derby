# Tasks: Scaffold del Proyecto

## Phase 1: Backend Foundation

- [x] 1.1 `go mod init` en `backend/` + `go.mod` con module path `github.com/ema/fixture`
- [x] 1.2 Crear estructura de directorios hexagonal (cmd, internal/{domain,handler,service,repository,router})
- [x] 1.3 Escribir tipos placeholder en `internal/domain/` (Categoria, Auto, Carrera)
- [x] 1.4 Definir interfaces de repositorio y servicio en `internal/repository/repository.go` y `internal/service/service.go`
- [x] 1.5 Agregar `github.com/go-chi/chi/v5` como dependencia
- [x] 1.6 Crear `internal/router/router.go` con router Chi configurado
- [x] 1.7 Crear `internal/handler/handler.go` (interfaz Handler) + `internal/handler/health.go` (health check)
- [x] 1.8 Crear `cmd/api/main.go` con servidor HTTP en puerto configurable (:8080 default)

## Phase 2: Frontend Foundation

- [x] 2.1 Crear proyecto SvelteKit 5 en `frontend/` con `npx sv create`
- [x] 2.2 Configurar Tailwind CSS en SvelteKit
- [x] 2.3 Habilitar PWA (vite-plugin-pwa)
- [x] 2.4 Configurar proxy de Vite hacia backend (:5173 → :8080)

## Phase 3: Dev Tooling

- [x] 3.1 Escribir `Makefile` con targets: `init`, `dev`, `build`, `test`, `vet`, `clean`, `all`
- [x] 3.2 Configurar Air (hot-reload) para Go backend
- [x] 3.3 Escribir `.gitignore` global
- [x] 3.4 Escribir `README.md` con instrucciones de setup

## Phase 4: Testing & Verify

- [x] 4.1 Test unitario para health endpoint (httptest) — 100% coverage
- [x] 4.2 Verificar `go vet ./...` pasa
- [x] 4.3 Verificar `gofmt -l .` sin errores
- [x] 4.4 Verificar `make dev` levanta backend + frontend
- [x] 4.5 Verificar `make build` produce binario

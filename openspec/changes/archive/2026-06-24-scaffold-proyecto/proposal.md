# Proposal: Scaffold del Proyecto

## Intent

Inicializar la estructura base del proyecto Go + Svelte 5 con arquitectura hexagonal, tooling de desarrollo, y configuración inicial de build para arrancar a codificar features reales.

**Req SRS asociados**: N/A (base habilitante para todos los RF)

## Scope

### In Scope
- Go module (`go mod init`) con estructura hexagonal (`cmd/`, `internal/`, `pkg/`)
- SvelteKit 5 frontend con PWA y Tailwind CSS
- Makefile con comandos: `dev`, `build`, `test`, `vet`, `clean`
- SQLite driver configurado en Go
- Config de dev: hot-reload (Air) + proxy inverso dev
- `.gitignore` + `README.md` inicial
- `go vet` + `gofmt` como quality gates

### Out of Scope
- Features de dominio (categorías, autos, carreras)
- WebSocket handlers
- UI/UX del frontend (solo scaffold)
- CI/CD (futuro cambio)
- Tests de integración o e2e (solo unit tests del scaffold)

## Capabilities

### New Capabilities
- `infra-base`: proyecto compilable con estructura hexagonal, tooling de dev y build funcionando

### Modified Capabilities
- None (primera capability del proyecto)

## Approach

Monorepo con dos módulos: `backend/` (Go) y `frontend/` (SvelteKit). El Makefile en la raíz orquesta ambos. Backend con arquitectura hexagonal: `cmd/api/main.go` como entrypoint, `internal/` con `handler/`, `service/`, `repository/`, `domain/`, `router/`. Frontend con SvelteKit + Tailwind + modo PWA habilitado.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/` | New | Go module completo |
| `frontend/` | New | SvelteKit project |
| `Makefile` | New | Comandos de dev/build/test |
| `.gitignore` | New | Ignorar bins, node_modules, .env |
| `README.md` | New | Setup instructions |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Svelte 5 breaking changes | Low | Usar create-svelte@latest con pinning |
| Go module path incorrecta | Low | Usar namespace claro |
| Puerto conflictivo dev | Low | Documentar y hacer configurable |

## Rollback Plan

`git init` → commit inicial → si falla, `rm -rf` y reiniciar. No hay datos ni código en juego.

## Dependencies

- Go 1.26.1+ (instalado ✅)
- Node 24.14.1+ (instalado ✅)
- npm (instalado ✅)

## Success Criteria

- [ ] `make dev` levanta backend y frontend simultáneamente
- [ ] `go vet ./...` pasa sin errores
- [ ] `go test ./...` pasa (aunque sean tests vacíos)
- [ ] Frontend SvelteKit devuelve página por defecto
- [ ] `make build` produce binario compilado

# Design: Scaffold del Proyecto

## Technical Approach

Monorepo con dos módulos independientes. Backend Go con arquitectura hexagonal (dominio puro, sin dependencias externas en `domain/`). Frontend SvelteKit 5 con Tailwind CSS y PWA configurado para mobile-first. Makefile como entrada única de desarrollo.

## Architecture Decisions

### Decision: Monorepo vs repos separados

| Opción | Tradeoff |
|--------|----------|
| **Monorepo** (elegido) | Makefile único, dev simplificado, un `git init` |
| Repos separados | CI/CD complejo, sincronización manual |
| **Rationale**: Proyecto chico, equipo de 1 persona. Monorepo reduce fricción.

### Decision: Estructura hexagonal vs MVC

| Opción | Tradeoff |
|--------|----------|
| **Hexagonal** (elegido) | `domain/` aislado, repositorios intercambiables, testeas sin infraestructura |
| MVC tradicional | Controller-Model-View, acoplamiento mayor |
| **Rationale**: SRS requiere persistencia ante cortes y portable — hexagonal permite cambiar SQLite por otro sin tocar dominio.

### Decision: SvelteKit con adapter-static vs adapter-node

| Opción | Tradeoff |
|--------|----------|
| **adapter-static** | Frontend 100% estático, sirve con Go embebido |
| adapter-node | Necesita Node runtime separado |
| **Rationale**: Go puede embeber el frontend compilado en el binario → un solo binario portable. Win-win.

## Data Flow (dev)

```
make dev
  ├── backend/ (Air hot-reload)
  │   └── :8080 (API)
  └── frontend/ (Vite dev)
      └── :5173 (proxy → :8080)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `backend/go.mod` | Create | Go module definition |
| `backend/cmd/api/main.go` | Create | Entrypoint HTTP server |
| `backend/internal/domain/carrera.go` | Create | Domain types placeholder |
| `backend/internal/router/router.go` | Create | Chi router setup |
| `backend/internal/handler/health.go` | Create | Health check endpoint |
| `backend/internal/handler/handler.go` | Create | Handler interface |
| `backend/internal/repository/repository.go` | Create | Repository interface |
| `backend/internal/service/service.go` | Create | Service interface |
| `frontend/` | Create | SvelteKit project (scaffolded) |
| `Makefile` | Create | Dev/build/test commands |
| `.gitignore` | Create | Ignore rules |
| `README.md` | Create | Setup instructions |

## Interfaces / Contracts

```go
// internal/domain/carrera.go
type Categoria struct {
    ID   string
    Nombre string
    EdadMin int
    EdadMax int
}

type Auto struct {
    ID       string
    Numero   int
    Nombre   string
    Creador  string
    Edad     int
    FotoURL  string
    CategoriaID string
}

type Carrera struct {
    ID         string
    CategoriaID string
    Autos      []string // IDs de autos
    OrdenLlegada []string // IDs ordenados
}
```

Esos tipos son **placeholder** — se refinarán en cambios posteriores. El scaffold solo define la estructura base.

## Testing Strategy

| Layer | Qué testear | Approach |
|-------|-------------|----------|
| Unit | Health endpoint | `httptest.NewRecorder` |
| Format | `gofmt -l` | Makefile gate |
| Lint | `go vet` | Makefile gate |

## Migration / Rollout

No migration required. Proyecto nuevo desde cero.

## Open Questions

- [None] Todo resuelto en decisiones de arquitectura.

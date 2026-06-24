# Design: Registro de Categorías

## Technical Approach

Backend: Agregar capa de persistencia SQLite con migraciones automáticas. Implementar el handler REST de categorías que usa un service concreto, que a su vez usa un repository SQLite. Frontend: Ruta `/admin/categorias` con listado y formulario modal.

## Architecture Decisions

### Decision: modernc.org/sqlite vs mattn/go-sqlite3

| Opción | Tradeoff |
|--------|----------|
| **modernc.org/sqlite** (elegido) | CGO-free → binario portable, build simple |
| mattn/go-sqlite3 | CGO required → build condicional |
| **Rationale**: SRS requiere portabilidad. `modernc.org/sqlite` es Go puro, produce un binario estático.

### Decision: Migraciones SQL vs ORM

| Opción | Tradeoff |
|--------|----------|
| **SQL plano con archivos** (elegido) | Simple, sin dependencias, control total |
| GORM / sqlx | Más dependencias, magia que complica debug |
| **Rationale**: Pocas tablas, migraciones simples. SQL plano alcanza.

### Decision: IDs con UUID vs autoincremental

| Opción | Tradeoff |
|--------|----------|
| **google/uuid** (elegido) | IDs únicos globales, seguros para API pública |
| autoincremental | Secuencial, predecible |
| **Rationale**: API REST y futura sincronización online. UUID evita colisiones.

## Data Flow

```
Cliente ──→ POST /api/categorias ──→ CategoriaHandler.Create ──→ CategoriaService.Create
                                                                       │
                                                              CategoriaRepository.Save
                                                                       │
                                                                   SQLite DB
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `backend/go.mod` | Modify | Agregar modernc.org/sqlite + google/uuid |
| `backend/internal/database/sqlite.go` | Create | Conexión + migraciones |
| `backend/internal/domain/carrera.go` | Modify | Agregar CreatedAt, UpdatedAt a Categoria |
| `backend/internal/repository/sqlite/categoria.go` | Create | Repository concreto SQLite |
| `backend/internal/repository/repository.go` | Modify | Agregar FindByID, Update, Delete |
| `backend/internal/service/categoria.go` | Create | Service concreto con validaciones |
| `backend/internal/service/service.go` | Modify | Agregar métodos faltantes al interface |
| `backend/internal/handler/categoria.go` | Create | Handler REST |
| `backend/cmd/api/main.go` | Modify | Inicializar DB, armar cadena repository→service→handler |
| `frontend/src/routes/admin/categorias/+page.svelte` | Create | Listado de categorías |
| `frontend/src/routes/admin/categorias/+page.server.ts` | Create | Load data SSR |
| `frontend/src/lib/components/CategoriaForm.svelte` | Create | Formulario creación/edición |

## Interfaces / Contracts

```go
// internal/database/sqlite.go
func New(dsn string) (*sql.DB, error)
func Migrate(db *sql.DB) error

// internal/repository/sqlite/categoria.go
type CategoriaRepository struct { db *sql.DB }
func (r *CategoriaRepository) List() ([]domain.Categoria, error)
func (r *CategoriaRepository) GetByID(id string) (*domain.Categoria, error)
func (r *CategoriaRepository) Save(c *domain.Categoria) error
func (r *CategoriaRepository) Update(c *domain.Categoria) error
func (r *CategoriaRepository) Delete(id string) error

// internal/service/categoria.go
type CategoriaService struct { repo domain.CategoriaRepository }
func (s *CategoriaService) List() ([]domain.Categoria, error)
func (s *CategoriaService) Create(nombre string, edadMin, edadMax int) (*domain.Categoria, error)
func (s *CategoriaService) Update(id string, nombre string, edadMin, edadMax int) (*domain.Categoria, error)
func (s *CategoriaService) Delete(id string) error
```

## Testing Strategy

| Layer | Qué testear | Approach |
|-------|-------------|----------|
| Unit | Handler | httptest con mocks de service |
| Unit | Service | Test con repo mock (slice en memoria) |
| Integration | Repository | SQLite en memoria (`:memory:`) |
| E2E | Frontend→Backend | Playwright (futuro) |

## Migration / Rollout

Migración inicial crea tabla `categorias`. Schema versionado con número secuencial en comentario SQL.

## Open Questions

- None resuelto en decisiones de arquitectura.

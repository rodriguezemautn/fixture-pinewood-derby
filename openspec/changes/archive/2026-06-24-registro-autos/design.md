# Design: Registro de Autos

## Technical Approach

Endpoint anidado: `/api/categorias/{categoriaId}/autos` para operaciones por categoría. Endpoint plano `/api/autos/{id}` para get/update/delete individual. Service valida existencia de categoría y unicidad de número. Foto almacenada como TEXT (URL/base64).

## Architecture Decisions

### Decision: Rutas anidadas vs planas

| Opción | Tradeoff |
|--------|----------|
| **Anidadas** (elegido) | `/api/categorias/{id}/autos` — semantic REST, consistente |
| Planas | `/api/autos?categoriaId=X` — menos RESTful |
| **Rationale**: Coherente con la relación jerárquica categoría→auto.

### Decision: Foto como TEXT vs archivos

| Opción | Tradeoff |
|--------|----------|
| **TEXT (URL/base64)** (elegido) | Simple, portable, sin file system |
| Archivos en disco | Complejidad, sync offline problemático |
| **Rationale**: SRS pide persistencia offline. Base64 en DB es más portable.

## Data Flow

```
POST /api/categorias/{id}/autos
  → AutoHandler.Create
    → verifica que categoría existe (CategoriaService.GetByID)
    → AutoService.Create
      → valida número único por categoría
      → AutoRepository.Save
        → INSERT INTO autos
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `backend/internal/database/sqlite.go` | Modify | Migración tabla autos |
| `backend/internal/domain/carrera.go` | Modify | Agregar CreatedAt, UpdatedAt a Auto |
| `backend/internal/repository/repository.go` | Modify | Agregar Update + ExistsByNumero a AutoRepository |
| `backend/internal/repository/sqlite/auto.go` | New | Repository SQLite |
| `backend/internal/service/auto.go` | New | Service con validaciones |
| `backend/internal/service/service.go` | Modify | Agregar interface AutoService |
| `backend/internal/handler/auto.go` | New | Handler REST |
| `backend/cmd/api/main.go` | Modify | Inicializar auto handler |
| `frontend/src/routes/admin/categorias/[id]/autos/` | New | Frontend pages |

## Interfaces / Contracts

```go
// AutoRepository extendido
type AutoRepository interface {
	ListByCategoria(categoriaID string) ([]domain.Auto, error)
	GetByID(id string) (*domain.Auto, error)
	Save(a *domain.Auto) error
	Update(a *domain.Auto) error
	Delete(id string) error
	ExistsByNumero(categoriaID string, numero int) (bool, error)
}

// AutoService
type AutoService interface {
	ListByCategoria(categoriaID string) ([]domain.Auto, error)
	GetByID(id string) (*domain.Auto, error)
	Create(categoriaID string, numero int, nombre, creador string, edad int, fotoURL string) (*domain.Auto, error)
	Update(id string, numero int, nombre, creador string, edad int, fotoURL string) (*domain.Auto, error)
	Delete(id string) error
}
```

## Testing Strategy

| Layer | Qué testear | Approach |
|-------|-------------|----------|
| Unit | Handler | httptest con mocks |
| Unit | Service | Mock repository |
| Integration | Repository | SQLite `:memory:` |

## Migration / Rollout

Nueva tabla `autos` con FK a categorías. No hay migración destructiva.

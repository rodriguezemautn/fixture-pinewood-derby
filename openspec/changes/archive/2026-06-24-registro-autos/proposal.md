# Proposal: Registro de Autos

## Intent

Implementar CRUD completo de autos asociados a categorías (RF2 del SRS), incluyendo subida de foto, backend completo y frontend de gestión.

**Req SRS**: RF2 — "El sistema deberá poder registrar uno o varios autos por categoría. Con su foto, número de auto, nombre de creador y edad."

## Scope

### In Scope
- Backend: handler REST para autos (CRUD por categoría)
- Backend: service con validaciones (número único por categoría, edad contra categoría)
- Backend: repository SQLite concreto con tabla autos
- Backend: migración DB para tabla autos + foreign key a categorías
- Backend: manejo de foto como URL/base64 (guardar en campo TEXT)
- Frontend: página de listado de autos por categoría
- Frontend: formulario de creación/edición con campos: número, nombre, creador, edad, foto
- Tests: handler, service, repository

### Out of Scope
- Subida real de archivos de imagen (se guarda URL/base64 por ahora)
- Scheduling de fixture
- Estilo gamer pulido (futuro cambio)
- Validación de edad contra categoría padre (validación simple por ahora)

## Capabilities

### New Capabilities
- `gestion-autos`: CRUD de autos con foto, asociados a categorías

### Modified Capabilities
- `gestion-categorias`: se agrega validación de integridad (no eliminar categoría con autos)

## Approach

Extender la arquitectura existente: nueva tabla `autos`, migración en database/sqlite.go, repositorio concreto, service con validaciones (número único por categoría), handler REST anidado bajo `/api/categorias/{categoriaId}/autos`. Frontend con ruta `/admin/categorias/{id}/autos`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/internal/database/sqlite.go` | Modified | Nueva migración tabla autos |
| `backend/internal/domain/carrera.go` | Modified | Agregar timestamps + CategoriaID validación |
| `backend/internal/repository/repository.go` | Modified | Agregar Update a AutoRepository |
| `backend/internal/repository/sqlite/auto.go` | New | Repository SQLite concreto |
| `backend/internal/service/auto.go` | New | Service con validaciones |
| `backend/internal/service/service.go` | Modified | Agregar interface AutoService |
| `backend/internal/handler/auto.go` | New | Handler REST anidado |
| `backend/cmd/api/main.go` | Modified | Inicializar auto repo/service/handler |
| `frontend/src/routes/admin/categorias/[id]/autos/` | New | Páginas frontend |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Foto en base64 pesada | Medium | Limitar tamaño, comprimir en frontend |
| Número auto duplicado en categoría | Low | Validación en service + unique constraint |

## Rollback Plan

`git revert` del commit. Si hay datos, borrar fixture.db.

## Dependencies

- `modernc.org/sqlite` ya instalado
- `google/uuid` ya instalado

## Success Criteria

- [ ] `GET /api/categorias/{id}/autos` lista autos de una categoría
- [ ] `POST /api/categorias/{id}/autos` crea auto con validaciones
- [ ] `GET /api/autos/{id}` retorna auto individual
- [ ] `PUT /api/autos/{id}` actualiza auto
- [ ] `DELETE /api/autos/{id}` elimina auto
- [ ] Validación: número único por categoría
- [ ] Frontend: listado + formulario con foto
- [ ] Tests > 70% coverage
- [ ] `go vet ./...` pasa

# Proposal: Registro de Categorías

## Intent

Implementar el CRUD completo de categorías (RF1 del SRS), incluyendo backend (handler, service, repository con SQLite) y frontend (listado, creación, edición). Base para todos los features posteriores.

**Req SRS asociado**: RF1 — "El sistema deberá poder registrar las categorías de la fecha."

## Scope

### In Scope
- Backend: handler HTTP con endpoints REST para categorías (CRUD)
- Backend: service con validaciones de negocio
- Backend: repository SQLite concreto (reemplaza la interfaz vacía)
- Backend: migración/inicialización de DB al arrancar
- Frontend: página de listado de categorías
- Frontend: formulario de creación/edición
- Frontend: diseño base "gamer/racing" con Tailwind
- Tests: unit tests para handler, service, repository

### Out of Scope
- Foto de categoría (no aplica)
- Scheduling de fixture (feature futuro)
- UI pulida con animaciones racing (se hará en cambio de styling)
- Autenticación de admin (futuro)

## Capabilities

### New Capabilities
- `gestion-categorias`: CRUD completo de categorías con persistencia SQLite

### Modified Capabilities
- `infra-base`: se modifica para agregar DB initialization + SQLite driver + migraciones

## Approach

Agregar SQLite con `modernc.org/sqlite` (CGO-free, portable). Repository concreto implementa las interfaces existentes. Handler RESTful con chi router. Frontend con SvelteKit: ruta `/admin/categorias` para gestión, ruta `/categorias` para vista pública.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/go.mod` | Modified | Agregar modernc.org/sqlite |
| `backend/cmd/api/main.go` | Modified | Inicializar DB, inyectar repository→service→handler |
| `backend/internal/domain/carrera.go` | Modified | Agregar timestamps a Categoria |
| `backend/internal/repository/repository.go` | Modified | CategoriaRepository concreta |
| `backend/internal/service/service.go` | Modified | CategoriaService concreta |
| `backend/internal/handler/` | New | categoria.go handler |
| `backend/internal/service/` | New | categoria.go service concreto |
| `backend/internal/repository/` | New | sqlite/categoria.go repositorio concreto |
| `backend/internal/database/` | New | sqlite.go conexión + migraciones |
| `frontend/src/routes/` | New | Paginas de gestión de categorías |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| modernc.org/sqlite incompatibilidad | Low | Usar versión estable, test rápido |
| Schema DB cambia frecuentemente | Medium | Migraciones con versionado simple |
| Frontend sin diseño racing | Low | Se iterará en cambio posterior |

## Rollback Plan

`git revert` del commit del cambio. Si hay datos en SQLite, borrar el archivo `.db` local.

## Dependencies

- `modernc.org/sqlite` (SQLite puro Go, sin CGO)
- `google/uuid` o similar para generación de IDs

## Success Criteria

- [ ] `POST /api/categorias` crea y retorna 201
- [ ] `GET /api/categorias` lista todas
- [ ] `GET /api/categorias/{id}` retorna una
- [ ] `PUT /api/categorias/{id}` actualiza
- [ ] `DELETE /api/categorias/{id}` elimina
- [ ] Validaciones: nombre requerido, edadMin < edadMax
- [ ] Frontend lista categorías y permite crear/editar
- [ ] Tests > 80% coverage en handler, service, repository
- [ ] `go vet ./...` pasa

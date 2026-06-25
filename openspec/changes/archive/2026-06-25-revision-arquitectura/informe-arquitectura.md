# Informe de Arquitectura — Revisión General

## Modelo de Datos (SQLite)

### Tablas

| Tabla | Propósito | Estado |
|-------|-----------|--------|
| `categorias` | Categorías de carrera por edad | ✅ |
| `autos` | Autos registrados por categoría | ✅ (con peso y foto) |
| `fixtures` | Fixtures de carrera (vinculados a competencias) | ⚠️ Fix aplicado (competencia_id) |
| `heats` | Carreras individuales dentro de un fixture | ✅ |
| `heat_autos` | Relación heats ↔ autos participantes | ✅ |
| `archivos_carrera` | Competencias archivadas | ✅ |
| `competencias` | Series de carreras por categoría | ✅ |

### Errores Corregidos

1. **BUG CRÍTICO**: `fixtures.competencia_id` nunca se seteaba
   - Síntoma: GetByCompetencia() nunca encontraba el fixture
   - Fix: Agregado `CompetenciaID` al domain, INSERT incluye columna, Create asigna valor
   - Commit: aa98c54

2. **BUG**: `GenerarFinal()` recibía categoriaId en vez de competenciaID
   - Síntoma: El endpoint no podía generar la carrera final
   - Fix: Movido a `/api/competencias/{id}/final`, frontend actualizado
   - Commit: aa98c54

3. **BUG**: Proxy Vite 8 bloqueaba requests POST/PUT/DELETE
   - Síntoma: Las mutaciones no llegaban al backend
   - Fix: Eliminado `changeOrigin: true` y `allowedHosts: true`
   - Commit: 96ec906

## API Endpoints

### Autenticación
```
POST /api/auth/login → 200 + token
```

### Categorías (CRUD)
```
GET    /api/categorias          → 200 []  (público)
POST   /api/categorias          → 201     (auth)
GET    /api/categorias/{id}     → 200     (público)
PUT    /api/categorias/{id}     → 200     (auth)
DELETE /api/categorias/{id}     → 204     (auth)
```

### Autos (CRUD)
```
GET    /api/autos               → 200 []  (público, todos)
GET    /api/categorias/{id}/autos → 200 [] (público, por categoría)
POST   /api/categorias/{id}/autos → 201   (auth)
GET    /api/autos/{id}          → 200     (público)
PUT    /api/autos/{id}          → 200     (auth)
DELETE /api/autos/{id}          → 204     (auth)
POST   /api/autos/{id}/foto     → 200     (auth, multipart)
```

### Competencias
```
GET    /api/categorias/{id}/competencias → 200 [] (público)
POST   /api/categorias/{id}/competencias?rondas=3 → 201 (auth)
POST   /api/competencias/{id}/finalizar  → 200 (auth)
POST   /api/competencias/{id}/desempate  → 201 (auth, body: auto_ids)
GET    /api/competencias/{id}/fixture    → 200 (público)
GET    /api/competencias/{id}/posiciones → 200 (público)
POST   /api/competencias/{id}/final      → 201 (auth)
```

### Fixture / Carreras
```
GET    /api/categorias/{id}/fixture      → 200 (público)
GET    /api/categorias/{id}/posiciones   → 200 (público)
POST   /api/carreras/{heatId}/resultado  → 200 (auth)
POST   /api/categorias/{id}/archivar     → 200 (auth)
GET    /api/categorias/{id}/archivos     → 200 (público)
```

### Salud
```
GET /health → 200
```

## Frontend Routes

| Ruta | Acceso | Propósito |
|------|--------|-----------|
| `/` | Público | Landing page con selector de carreras |
| `/login` | Público | Login admin |
| `/admin/categorias` | Auth | CRUD categorías |
| `/admin/categorias/{id}/autos` | Auth | CRUD autos + selector competencias |
| `/admin/categorias/{id}/fixture` | Auth | Fixture + resultados + desempate |
| `/carreras` | Público | Lista de categorías con carreras |
| `/carreras/{id}` | Público | Vista pública de carrera (heats, posiciones, podio) |
| `/autos` | Público | Todos los autos registrados (agrupados por categoría) |

## Flujo de Actividades

```
Categoría
  └── Registrar Autos (mín 4)
        └── Crear Competencia #1
              ├── Generar Heats (3 rondas, shuffle)
              ├── Registrar Resultados de cada Heat
              ├── Generar Carrera Final
              ├── Registrar Resultado Final
              ├── Finalizar Competencia
              │     └── (opcional) Agregar Desempate
              └── Crear Competencia #2 (con distintos grupos)
```

## Tests

| Paquete | Tests |
|---------|-------|
| auth | 1 |
| handler | 15 |
| repository/sqlite | 9 |
| service | 23 (9 scheduler + 14 service) |
| **Total** | **48** |

# Documento de Arquitectura — Fixture D15

> **Versión**: 1.0  
> **Última actualización**: 2026-06-25  
> **Stack**: Go 1.26 + Svelte 5 + SQLite + Tailwind CSS 4

---

## 1. Stack Tecnológico

| Capa | Tecnología | Versión | Propósito |
|------|-----------|---------|-----------|
| Backend | Go | 1.26.1 | API REST, lógica de negocio, scheduler |
| Frontend | Svelte | 5.x | UI reactiva con runes mode |
| CSS | Tailwind | 4.x | Utility-first CSS + diseño propio |
| Base de datos | SQLite (modernc.org) | — | Persistencia embebida, CGO_ENABLED=0 |
| Router HTTP | chi | v5 | Enrutamiento REST |
| Auth | HMAC-SHA256 | — | Tokens expirables para admin |
| PWA | vite-plugin-pwa | — | Service worker, offline support |
| Build | Vite | 8.x | Bundler frontend |

### Requisitos del sistema

- **Go** 1.22+ (compilación estática con CGO_ENABLED=0)
- **Node.js** 20+ (para frontend en desarrollo)
- **npm** 10+
- **sqlite3** CLI (opcional, para consultas directas)
- **firewalld** (opcional, para abrir puertos en red)

---

## 2. Estructura del Proyecto

```
/
├── backend/
│   ├── cmd/api/main.go           ← Entrypoint del servidor HTTP
│   ├── internal/
│   │   ├── auth/                 ← Middleware JWT-like HMAC
│   │   ├── database/             ← Conexión SQLite + migraciones
│   │   ├── domain/               ← Structs de dominio
│   │   ├── handler/              ← Handlers HTTP (REST)
│   │   ├── repository/
│   │   │   ├── repository.go     ← Interfaces de repositorio
│   │   │   └── sqlite/           ← Implementaciones SQLite
│   │   ├── router/               ← Router chi con middleware
│   │   └── service/              ← Lógica de negocio
│   └── tmp/main                  ← Binario compilado
│
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api.ts            ← Helper fetch con auth
│   │   │   ├── components/       ← Componentes reutilizables
│   │   │   └── assets/           ← SVG, favicon
│   │   └── routes/               ← Páginas (SvelteKit file-based)
│   ├── static/assets/            ← Imágenes
│   └── vite.config.ts
│
├── manager/
│   ├── main.go                   ← Consola web de gestión (independiente)
│   └── dashboard.html            ← Dashboard HTML embebido
│
├── docs/                         ← Documentación
├── bin/manage.sh                 ← Script de gestión por terminal
├── SRS.md                        ← Especificación de Requisitos
└── Makefile                      ← Comandos de desarrollo
```

---

## 3. Modelo de Datos

### Diagrama Entidad-Relación (ASCII)

```
┌─────────────┐       ┌─────────────┐       ┌──────────────────┐
│  categorias │───────│    autos    │       │  archivos_carrera │
│─────────────│       │─────────────│       │──────────────────│
│ id (PK)     │       │ id (PK)     │       │ id (PK)          │
│ nombre      │       │ categoria_id│◄──────│ categoria_id     │
│ edad_min    │       │ numero      │       │ categoria_nombre │
│ edad_max    │       │ nombre      │       │ fecha            │
│ created_at  │       │ creador     │       │ winner_id        │
│ updated_at  │       │ edad        │       │ winner_nombre    │
└──────┬──────┘       │ peso        │       │ winner_numero    │
       │              │ foto_url    │       │ resultados (JSON)│
       │              │ created_at  │       └──────────────────┘
       │              │ updated_at  │
       │              └──────┬──────┘
       │                     │
┌──────┴──────────────┐      │
│    competencias     │      │
│─────────────────────│      │
│ id (PK)             │      │
│ categoria_id ───────┘      │
│ numero                    │
│ nombre                    │
│ estado (abierta/fin)      │
│ rondas                    │
│ created_at                │
└──────┬───────────────────┘
       │
┌──────┴───────────────────┐   ┌──────────────────┐
│       fixtures           │   │      heats       │
│──────────────────────────│   │──────────────────│
│ id (PK)                  │   │ id (PK)          │
│ categoria_id ────────────┘   │ fixture_id ──────┘
│ competencia_id ──────────┐   │ numero           │
│ rondas                   │   │ completado (bool)│
│ estado (pendiente/fin)   │   │ orden_llegada    │
│ created_at               │   │ registrado_at    │
└──────────────────────────┘   │ created_at       │
                               └────────┬─────────┘
                                        │
                               ┌────────┴─────────┐
                               │    heat_autos    │
                               │──────────────────│
                               │ heat_id (PK)     │
                               │ auto_id (PK) ────┘
                               └──────────────────┘
```

### Relaciones Clave

| Relación | Tipo | Descripción |
|----------|------|-------------|
| categorias → autos | 1:N | Una categoría tiene muchos autos |
| categorias → competencias | 1:N | Una categoría tiene muchas competencias |
| categorias → archivos_carrera | 1:N | Historial de carreras archivadas |
| competencias → fixtures | 1:1 | Una competencia tiene exactamente un fixture |
| fixtures → heats | 1:N | Un fixture tiene muchos heats (carreras) |
| heats → heat_autos → autos | M:N | Muchos autos participan en muchos heats |

### Migraciones

Las migraciones son **idempotentes** y se ejecutan al iniciar el backend.
Se usa `CREATE TABLE IF NOT EXISTS` y `ALTER TABLE` con skip de errores.

```go
// Ejemplo de migración
`CREATE TABLE IF NOT EXISTS competencias (
    id TEXT PRIMARY KEY,
    categoria_id TEXT NOT NULL,
    numero INTEGER NOT NULL DEFAULT 1,
    nombre TEXT NOT NULL DEFAULT '',
    estado TEXT NOT NULL DEFAULT 'abierta',
    ...
    FOREIGN KEY (categoria_id) REFERENCES categorias(id)
)`,
```

> **Nota**: Desde el commit 67b5e5e, `PRAGMA foreign_keys = ON` está habilitado.

---

## 4. API REST

### Autenticación

```
POST /api/auth/login
  Body: { "pin": "1234" }
  → 200 { "token": "base64...", "role": "admin" }
```

El token HMAC tiene validez de **24 horas**. Los endpoints de escritura (POST/PUT/DELETE) requieren `Authorization: Bearer <token>`. Los GET son públicos.

### Endpoints

#### Categorías
```
GET    /api/categorias            → 200 []    (público)
POST   /api/categorias            → 201       (auth)
GET    /api/categorias/{id}       → 200       (público)
PUT    /api/categorias/{id}       → 200       (auth)
DELETE /api/categorias/{id}       → 204       (auth, protegido)
```

#### Autos
```
GET    /api/autos                 → 200 []    (todos)
GET    /api/categorias/{id}/autos → 200 []    (por categoría)
POST   /api/categorias/{id}/autos → 201       (auth)
GET    /api/autos/{id}            → 200       (público)
PUT    /api/autos/{id}            → 200       (auth)
DELETE /api/autos/{id}            → 204       (auth, protegido)
POST   /api/autos/{id}/foto       → 200       (auth, multipart)
```

#### Competencias
```
GET    /api/categorias/{id}/competencias → 200 []    (público)
POST   /api/categorias/{id}/competencias → 201       (auth, ?rondas=N)
POST   /api/competencias/{id}/finalizar  → 200       (auth)
POST   /api/competencias/{id}/desempate  → 201       (auth)
GET    /api/competencias/{id}/fixture    → 200       (público)
GET    /api/competencias/{id}/posiciones → 200       (público)
POST   /api/competencias/{id}/final      → 201       (auth)
```

#### Fixture / Carreras
```
GET    /api/categorias/{id}/fixture     → 200       (público)
GET    /api/categorias/{id}/posiciones  → 200       (público)
POST   /api/carreras/{heatId}/resultado → 200       (auth)
POST   /api/categorias/{id}/archivar    → 200       (auth, body: competencia_id)
GET    /api/categorias/{id}/archivos    → 200       (público)
```

#### Salud
```
GET /health → 200 { "status": "ok" }
```

---

## 5. Frontend Routes

| Ruta | Acceso | Descripción |
|------|--------|-------------|
| `/` | Público | Landing page con escenario arcade racing |
| `/login` | Público | Login con PIN |
| `/admin/categorias` | Auth | CRUD de categorías |
| `/admin/categorias/{id}/autos` | Auth | CRUD de autos + selector de competencias |
| `/admin/categorias/{id}/fixture` | Auth | Fixture, resultados, desempate, podio |
| `/carreras` | Público | Lista de categorías disponibles |
| `/carreras/{id}` | Público | Vista pública: heats, posiciones, podio, celebración |
| `/autos` | Público | Todos los autos agrupados por categoría |

---

## 6. Flujo de Autenticación

```
Usuario → /login → ingresa PIN
  → POST /api/auth/login
    → backend valida PIN contra env o "1234"
    → genera token HMAC(admin:timestamp, secret)
    → devuelve { token, role }
  → frontend guarda en localStorage("auth_token")
  → redirige a /admin/categorias

En cada request POST/PUT/DELETE:
  → apiFetch() lee token de localStorage
  → agrega header "Authorization: Bearer <token>"
  → backend verifica HMAC + expiración

Si el token expiró:
  → backend responde 401
  → apiFetch() borra token y redirige a /login
```

---

## 7. Algoritmo de Fixture (Swiss-System)

El fixture se genera con un algoritmo Swiss-system adaptado:

```
Entrada: lista de autoIDs, rondas R
Salida: R * ceil(N/4) heats de 4 autos

1. Si hay menos de 4 autos, error
2. Mezclar autos aleatoriamente (primera ronda)
3. Para cada ronda:
   a. Ordenar autos por puntos ascendente
   b. Agrupar en bloques de 4 (head-to-head)
   c. Para el último bloque, distribuir autos sobrantes
4. Asignar carriles aleatoriamente dentro de cada heat
```

**Puntuación**: 1°=5pts, 2°=3pts, 3°=2pts, 4°=1pt, DNS=0pts
**Desempate**: Mayor cantidad de 1° puestos → menor edad

---

## 8. Tests

| Paquete | Cantidad | Cobertura |
|---------|:--------:|:---------:|
| auth | 1 | Middleware, generación/verificación de tokens |
| handler | 15 | CRUD categorías, autos, fixture |
| repository/sqlite | 9 | Persistencia SQLite |
| service | 23 | Lógica de negocio, scheduler, cálculos |
| **Total** | **48** | Todos verdes |

Ejecutar: `make test` o `cd backend && go test -v ./...`

---

## 9. Seguridad e Integridad

- **Foreign keys**: `PRAGMA foreign_keys = ON` desde commit 67b5e5e
- **Borrado protegido**: No se puede eliminar una categoría con autos, competencias o archivos
- **Borrado de autos**: No se puede eliminar un auto que fue campeón en una carrera archivada
- **Token expirable**: 24h de validez, redirección automática al expirar
- **CGO_ENABLED=0**: Binario 100% estático, sin dependencias del sistema
- **SQLite WAL mode**: Mejor concurrencia en lecturas concurrentes

---

## 10. Consola de Gestión

El proyecto incluye dos herramientas de gestión independientes:

| Herramienta | Tipo | Puerto | Descripción |
|-------------|------|:------:|-------------|
| `bin/manage.sh` | Script bash | — | Menú interactivo por terminal |
| `manager/` | Web app Go | 9099 | Dashboard web con logs en vivo |

Ambas permiten: iniciar/detener servicios, ver estado, monitorear logs, compilar, resetear DB y abrir puertos firewall.

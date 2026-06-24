# 🏎️ Fixture — Pinewood Derby Racing

Controlador de fixture para carreras de autitos de madera (pinewood derby) del **Destacamento 15 de la Iglesia Betel**.

> Sistema completo de gestión: categorías, autos, scheduler de carreras, podios animados y vista pública para proyección.

---

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | **Go 1.26** — Arquitectura Hexagonal |
| Frontend | **Svelte 5** + SvelteKit 2 + Tailwind CSS |
| PWA | vite-plugin-pwa (service worker + manifest) |
| BD | **SQLite** via modernc.org/sqlite (CGO-free) |
| Hot-reload | Air |
| Calidad | go vet + gofmt + **58 tests** |

## Requisitos

- Go 1.26+
- Node 24+
- npm

## Inicio rápido

```bash
# 1. Clonar e inicializar
git clone <repo> && cd fixture
make init          # go mod tidy + npm install

# 2. Iniciar backend (Terminal 1)
make dev-backend   # Air hot-reload → http://localhost:8080

# 3. Iniciar frontend (Terminal 2)
make dev-frontend  # Vite dev server → http://localhost:5173
```

### Si el puerto 8080 está ocupado

```bash
make kill-backend  # Mata proceso en :8080 automáticamente
make dev-backend   # Arranca de nuevo
```

### Build + tests

```bash
make build    # Compila backend + frontend
make test     # go test -v -cover ./backend/...
make vet      # go vet ./...
make fmt      # gofmt -l .
make all      # Pipeline completo: clean → vet → test → build
```

## Credenciales

| Rol | PIN | Login |
|-----|-----|-------|
| Admin | `1234` | http://localhost:5173/login |

Configurable vía variable de entorno:

```bash
ADMIN_PIN=5678 make dev-backend
```

## Rutas del frontend

| Ruta | Auth | Descripción |
|------|------|-------------|
| `/` | No | Landing racing con hero e info cards |
| `/login` | No | Login de administrador |
| `/admin/categorias` | Sí* | Gestión de categorías |
| `/admin/categorias/{id}/autos` | Sí* | Autos por categoría |
| `/admin/categorias/{id}/fixture` | Sí* | Fixture + tabla de posiciones |
| `/carreras/{id}` | No | **Vista pública** para proyección 🏁 |

*\*El frontend redirige a `/login` si no hay token. Las rutas GET de la API son públicas; solo POST/PUT/DELETE requieren autenticación.*

## API endpoints

```
Autenticación
  POST /api/auth/login        → 200 + token (body: {"pin":"1234"})

Categorías (CRUD)
  GET    /api/categorias              → 200 []       (público)
  POST   /api/categorias              → 201          (auth)
  GET    /api/categorias/{id}         → 200          (público)
  PUT    /api/categorias/{id}         → 200          (auth)
  DELETE /api/categorias/{id}         → 204          (auth)

Autos (CRUD)
  GET    /api/categorias/{id}/autos   → 200 []       (público)
  POST   /api/categorias/{id}/autos   → 201          (auth)
  GET    /api/autos/{id}              → 200          (público)
  PUT    /api/autos/{id}              → 200          (auth)
  DELETE /api/autos/{id}              → 204          (auth)

Fixture
  POST   /api/categorias/{id}/fixture?rondas=3  → 201 (auth)
  GET    /api/categorias/{id}/fixture           → 200 (público)
  GET    /api/categorias/{id}/posiciones        → 200 (público)
  POST   /api/carreras/{heatId}/resultado       → 200 (auth)
  POST   /api/categorias/{id}/final             → 201 (auth)

Salud
  GET    /health                        → 200
```

## Arquitectura del backend

```
cmd/api/main.go                ← Entrypoint: inicializa DB, inyecta dependencias
└── internal/
    ├── auth/                   ← Middleware HMAC + generación de tokens
    ├── database/               ← Conexión SQLite + migraciones automáticas
    ├── domain/                 ← Entidades puras (Categoria, Auto, Carrera, Fixture)
    ├── handler/                ← Handlers HTTP (REST endpoints)
    ├── repository/             ← Interfaces de persistencia
    │   └── sqlite/             ← Implementaciones concretas SQLite
    ├── router/                 ← Chi router + middleware global
    └── service/                ← Lógica de negocio + algoritmo scheduler
```

### Algoritmo de fixture (Swiss-system)

```
Rondas: R (= 3 default)
Heats: 4 autos por heat
Puntos: 1°=5, 2°=3, 3°=2, 4°=1, DNS=0
Emparejamiento: Swiss-system (autos de similar performance compiten juntos)
Desempate: más 1° puestos → menor edad
Top 4 → Carrera Final 🏆
```

## Estructura del proyecto

```
fixture/
├── backend/               ← Go 1.26 (hexagonal)
│   ├── cmd/api/main.go
│   └── internal/{auth,database,domain,handler,repository,router,service}
├── frontend/              ← SvelteKit 2 + Tailwind + PWA
│   └── src/
│       ├── lib/components/   ← Podium, Celebration, CategoriaForm, AutoForm
│       ├── lib/api.ts        ← Fetch helper con auth
│       └── routes/           ← Landing, Admin, Login, Carreras
├── openspec/              ← SDD artifacts
│   ├── config.yaml
│   ├── specs/             ← Specs principales
│   └── changes/archive/   ← 7 cambios completados
├── public/assets/         ← Logos D15, emblemas, autos de madera
├── Makefile               ← init dev build test vet fmt clean all
└── SRS.md                 ← Requerimientos funcionales originales
```

## Tests

```bash
# Todos los tests
go test ./backend/... -v -count=1    # 58 tests

# Por paquete
go test ./backend/internal/handler/...
go test ./backend/internal/service/...
go test ./backend/internal/repository/sqlite/...
go test ./backend/internal/auth/...
```

| Paquete | Tests | Cobertura |
|---------|-------|-----------|
| handler | 15 | ~50% |
| service | 22 | ~66% |
| repository | 9 | ~31% |
| auth | 1 | — |

## Funcionalidades por SRS

| Req | Descripción | Estado |
|-----|-------------|--------|
| RF1 | Registrar categorías | ✅ |
| RF2 | Registrar autos con foto | ✅ |
| RF3 | Máximo 4 autos por carrera | ✅ |
| RF4 | Registrar orden de llegada | ✅ |
| RF5 | Scheduler fixture óptimo y equitativo | ✅ |
| RF6 | Mostrar fixture gráfico | ✅ |
| RF7 | Podio emocionante | ✅ |
| RF8 | Ventana especial podio final | ✅ |
| RB1 | Admin gestiona categorías y carreras | ✅ |
| RB2 | Visualizador solo online | ✅ |

## Créditos

Desarrollado para el **Destacamento 15** de la Iglesia Betel. Pinewood Derby — una tradición de exploradores del Rey.

Assets y logos propiedad del Destacamento 15.

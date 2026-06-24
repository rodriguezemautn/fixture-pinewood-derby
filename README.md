# 🏎️ Fixture — Pinewood Derby Racing

Controlador de fixture para carreras de autitos derby del Destacamento 15 de la Iglesia Betel.

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go 1.26+ (Hexagonal Architecture) |
| Frontend | Svelte 5 + SvelteKit 2 + Tailwind CSS |
| PWA | vite-plugin-pwa (offline support) |
| BD | SQLite (vía CGO) |
| Dev | Air (hot-reload) |

## Requisitos

- Go 1.26+
- Node 24+
- npm

## Desarrollo

```bash
# Inicializar proyecto (solo primera vez)
make init

# Levantar backend + frontend simultáneamente
make dev

# Compilar todo
make build

# Correr tests
make test

# Linter
make vet

# Verificar formato
make fmt

# Limpiar artefactos
make clean

# Build + test completo
make all
```

Backend corre en `:8080`, frontend en `:5173` con proxy a `:8080`.

## Estructura

```
backend/
├── cmd/api/           ← Entrypoint HTTP
├── internal/
│   ├── domain/        ← Entidades del negocio
│   ├── handler/       ← Handlers HTTP
│   ├── repository/    ← Interfaces de persistencia
│   ├── router/        ← Configuración de rutas
│   └── service/       ← Lógica de negocio
└── pkg/               ← Utilidades exportables

frontend/
├── src/
│   ├── routes/        ← Páginas SvelteKit
│   └── lib/           ← Componentes reutilizables
└── build/             ← Frontend compilado
```

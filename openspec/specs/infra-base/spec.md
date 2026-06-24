# Infra-base Specification

## Purpose

Define la estructura base del proyecto, tooling de desarrollo, y configuración de build que habilita la implementación de features de dominio.

## Requirements

### Requirement: Estructura de directorios

El proyecto **MUST** organizar el backend en estructura hexagonal: `cmd/`, `internal/{handler,service,repository,domain,router}`, `pkg/`.

#### Scenario: Backend structure created

- GIVEN un proyecto Go vacío
- WHEN se ejecuta `make init`
- THEN los directorios `backend/cmd/api/`, `backend/internal/handler/`, `backend/internal/service/`, `backend/internal/repository/`, `backend/internal/domain/`, `backend/internal/router/` existen
- AND `backend/go.mod` existe con module path definido

### Requirement: SvelteKit frontend

El proyecto **MUST** incluir un frontend SvelteKit 5 con Tailwind CSS y configuración PWA.

#### Scenario: Frontend scaffolded

- GIVEN Node 24+ instalado
- WHEN se ejecuta `make init`
- THEN `frontend/` contiene un proyecto SvelteKit funcional
- AND `frontend/package.json` existe con dependencias de Svelte 5 y Tailwind
- AND `npm run dev` inicia el dev server sin errores

### Requirement: Makefile con comandos dev

El proyecto **SHALL** tener un Makefile con comandos para desarrollo y build.

#### Scenario: Dev command works

- GIVEN el scaffold completo
- WHEN se ejecuta `make dev`
- THEN backend y frontend se inician simultáneamente

#### Scenario: Build command works

- GIVEN el scaffold completo
- WHEN se ejecuta `make build`
- THEN se produce un binario en `backend/bin/`

### Requirement: Quality gates

El proyecto **MUST** pasar `go vet` y `gofmt` sin errores.

#### Scenario: Vet passes

- GIVEN código Go compilable
- WHEN se ejecuta `make vet`
- THEN `go vet ./...` retorna exit code 0

#### Scenario: Format check passes

- GIVEN código Go formateado
- WHEN se ejecuta `gofmt -l .`
- THEN no retorna archivos sin formatear

### Requirement: Inicialización de base de datos

El proyecto **MUST** inicializar una base de datos SQLite al arrancar el servidor, ejecutando migraciones automáticas para crear las tablas necesarias.

#### Scenario: DB se inicializa al iniciar

- GIVEN el servidor arranca por primera vez
- WHEN se ejecuta `cmd/api/main.go`
- THEN se crea el archivo `fixture.db` en la raíz del proyecto
- AND la tabla `categorias` existe con columnas: id, nombre, edad_min, edad_max, created_at, updated_at

#### Scenario: Migraciones idempotentes

- GIVEN la DB ya existe con tablas creadas
- WHEN el servidor se reinicia
- THEN no hay errores
- AND los datos existentes se conservan

### Requirement: SQLite driver sin CGO

El proyecto **MUST** usar `modernc.org/sqlite` como driver SQLite para mantener portabilidad sin dependencia de CGO.

#### Scenario: Build sin CGO

- GIVEN el entorno sin GCC/CGO
- WHEN se ejecuta `CGO_ENABLED=0 go build`
- THEN el binario se compila sin errores

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

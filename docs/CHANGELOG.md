# Changelog — Fixture D15

> Todos los commits del proyecto, agrupados por funcionalidad.
> Formato: [Conventional Commits](https://conventionalcommits.org/)

---

## [1.0.0] — 2026-06-25 — Sprint Inicial

### 🏗️ Scaffold y Configuración Inicial

| Commit | Descripción |
|--------|-------------|
| `1274652` | Experiencia racing en vivo (polling, skeletons, breadcrumbs) |
| `90ab4bc` | Peso del auto en gramos + fix /autos endpoint |
| `6229410` | Fix error al archivar competencia |

### 📋 Registro de Categorías

| Commit | Descripción |
|--------|-------------|
| *(scaffold)* | CRUD categorías con persistencia SQLite + frontend admin |
| *(specs)* | Specs infra-base y gestión-categorías |

### 🚗 Registro de Autos

| Commit | Descripción |
|--------|-------------|
| *(scaffold)* | CRUD autos asociados a categorías con foto |
| *(specs)* | Spec gestión-autos |

### 🎨 Estilo Racing

| Commit | Descripción |
|--------|-------------|
| `a7c0162` | **Frontend restyle arcade racing 8-bit** — paleta naranja retro, tipografía Press Start 2P + VT323, scanlines CRT, bordes pixel |
| `f9ec993` | **Tema claro (naranja) + oscuro mejorado** + escenario arcade racing con estación Berazategui |
| `fe41034` | Ajuste de paleta para proyector (fondos más claros, texto con más contraste) |
| `0a1dddd` | Título blanco con sombra para contraste, estación mejorada, cartel BERAZATEGUI negro/blanco, tren azul animado |
| `5e7f050` | Fábrica con chimenea blanca/roja, tren bala con vagones que para en estación |
| `fc3ea54` | Fábrica más grande con paredes naranjas y puerta |
| `37996d7` | Título "Autitos Derby Exploradores del Rey", Destacamento 15 más grande |
| `e9d30df` | Fix contraste: todas las líneas del título en blanco con sombra |
| `abdc85c` | Categoría y fecha en podio + celebración "Campeón de Categoría" |

### 🏁 Sistema de Competencias

| Commit | Descripción |
|--------|-------------|
| `9ad702a` | Timestamp en heats + archivar competencia |
| `fcbd1ed` | **Sistema de competencias por categoría** |
| `853b73f` | Fix algoritmo fixture con rondas variadas + finalWinner correcto |

### 🐛 Fixes de Bugs

| Commit | Descripción |
|--------|-------------|
| `aa98c54` | **CRÍTICO**: `fixture.competencia_id` nunca se seteaba + GenerarFinal con ID incorrecto |
| `96ec906` | Proxy Vite 8 sin changeOrigin (bloqueaba POST/PUT/DELETE) |
| `4091ec6` | Feedback visual en nueva competencia + error handling |
| `8e8b400` | **Archivar ahora acepta competencia_id** y limpia solo ese fixture |
| `c8fce11` | Redirección a login cuando el token expira (401 en apiFetch) |
| `fa11da4` | Modales con desborde vertical (max-height 85vh + overflow-y auto) |
| `afafba6` | Botón Celebrar en vista pública muestra categoría y fecha |

### 🚄 Desempate y Podio

| Commit | Descripción |
|--------|-------------|
| `138ec25` | **Permitir desempate sin finalizar** + botón Ver Podio Final con celebración |

### 🔐 Integridad y Seguridad

| Commit | Descripción |
|--------|-------------|
| `67b5e5e` | **PRAGMA foreign_keys ON** + protección borrado categorías/autos con historial |

### 🛠️ Gestión y Monitoreo

| Commit | Descripción |
|--------|-------------|
| `f8c3728` | Script de gestión `bin/manage.sh` (menú interactivo) |
| `2ca5333` | **Consola web de gestión independiente** (`manager/`) — dashboard con logs SSE |
| `81344ea` | Dos paneles de logs en paralelo (backend + frontend) + botones deshabilitados según estado |

### 📚 Documentación

| Commit | Descripción |
|--------|-------------|
| `a2ac1be` | Informe de arquitectura, tareas pendientes y avance |
| `b5ca06b` | Organización de participantes por categoría |
| `923cb6f` | Campeón de Campeones vacío hasta definir ganadores |
| `2c9352e` | Eliminar Gonzalez Jennifer (protección de datos) |
| `0921176` | Columnas Iglesia (Betel) y Localidad (Berazategui) en planillas |
| *(este commit)* | Documentos de arquitectura, diseño y técnico |

---

## Estadísticas

| Métrica | Valor |
|---------|:-----:|
| **Commits totales** | 53 |
| **Tests** | 48 (todos verdes) |
| **Backend** | Go 1.26, ~35 archivos |
| **Frontend** | Svelte 5, ~20 archivos |
| **Base de datos** | SQLite, 7 tablas |
| **Binario** | Estático, ~16MB |

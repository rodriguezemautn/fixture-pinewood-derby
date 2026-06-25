# Informe de Avance — Sprint Final v1.0

**Fecha**: 2026-06-25  
**Commits**: 53 (desde inicio)

---

## Resumen Ejecutivo

Sistema completo de gestión de carreras pinewood derby para el Destacamento 15.
Se implementaron todas las funcionalidades requeridas, más herramientas de gestión
y documentación técnica. El sistema está listo para uso en producción.

---

## Funcionalidades Implementadas

### Core
- ✅ **Registro de categorías** por edad (5 rangos + Campeón de Campeones)
- ✅ **Registro de autos** con foto, peso, número único por categoría
- ✅ **Sistema de competencias múltiples** por categoría
- ✅ **Fixture Swiss-system** con R rondas, heats de 4 autos
- ✅ **Registro de resultados** con orden de llegada
- ✅ **Carrera final** con top 4
- ✅ **Podio animado** con celebración (🥇🥈🥉)
- ✅ **Desempate** antes o después de finalizar
- ✅ **Archivar competencia** con historial
- ✅ **Autenticación** por PIN con token HMAC
- ✅ **Vista pública** con polling en vivo cada 10s

### UI/UX
- ✅ **Tema oscuro** arcade racing con alto contraste
- ✅ **Tema claro** naranja retro (proyector-friendly)
- ✅ **Escenario arcade CSS** en landing (ruta, fábrica, estación, tren)
- ✅ **Scanlines CRT**, bordes pixel, glow, animaciones retro
- ✅ **Responsive** mobile + proyector

### Gestión y Monitoreo
- ✅ **Script por terminal** (`bin/manage.sh`) con menú interactivo
- ✅ **Consola web independiente** (`manager/`) en puerto 9099
- ✅ Logs en vivo (backend + frontend en paralelo)
- ✅ Botones inteligentes (Start/Stop según estado)
- ✅ Firewall, compilación, reset DB desde la UI

### Integridad
- ✅ `PRAGMA foreign_keys = ON`
- ✅ Protección de borrado de categorías con historial
- ✅ Protección de borrado de autos campeones
- ✅ Redirección automática al login cuando el token expira

---

## Métricas

| Métrica | Anterior | Actual |
|---------|:--------:|:------:|
| Commits | 9 | **53** |
| Tests | 48 | **48** (todos verdes) |
| Backend (archivos Go) | ~29 | **~35** |
| Frontend (componentes) | ~17 | **~20** |
| Líneas de backend | — | **~3500** |
| Líneas de frontend | — | **~5000** |
| Documentación | — | **8 archivos** |
| Binario backend | — | **~16MB estático** |

---

## Bugs Corregidos (post-arquitectura)

| Bug | Fix | Commit |
|-----|-----|--------|
| Modal sin scroll al cargar foto | max-height + overflow-y | `fa11da4` |
| Token expirado no redirige al login | Detección de 401 en apiFetch | `c8fce11` |
| Archivar ignoraba competencia_id | Parsear body en handler | `8e8b400` |
| Desempate solo después de finalizar | Eliminar restricción de estado | `138ec25` |
| Borrar categoría con datos | HasArchivos/HasCompetencias/HasAutos | `67b5e5e` |
| Borrar auto campeón | IsWinnerInArchive | `67b5e5e` |

---

## Cobertura SRS

| RF | Descripción | Estado |
|----|-------------|--------|
| RF1 | Registrar categorías | ✅ |
| RF2 | Registrar autos con foto | ✅ |
| RF3 | Máximo 4 autos por carrera | ✅ |
| RF4 | Registrar orden de llegada | ✅ |
| RF5 | Scheduler fixture óptimo | ✅ Swiss-system |
| RF6 | Mostrar fixture gráfico | ✅ |
| RF7 | Podio emocionante | ✅ |
| RF8 | Ventana especial podio final | ✅ |
| RB1 | Admin gestiona | ✅ auth PIN |
| RB2 | Visualizador online | ✅ vista pública |

---

## Documentación Generada

| Documento | Contenido |
|-----------|-----------|
| `docs/arquitectura.md` | Stack, modelo de datos, API, frontend routes, flujos |
| `docs/diseno.md` | Sistema de diseño, colores, tipografía, componentes |
| `docs/tecnico.md` | Instalación, desarrollo, build, troubleshooting |
| `docs/CHANGELOG.md` | Todos los commits agrupados por funcionalidad |
| `SRS.md` | Especificación de Requisitos del Sistema |

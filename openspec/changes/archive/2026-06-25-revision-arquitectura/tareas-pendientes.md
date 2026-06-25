# Tareas Pendientes — Próximos Pasos

**Actualizado**: 2026-06-25

---

## Estado General

La funcionalidad core está **completa**. Las tareas pendientes son
principalmente deuda técnica y mejoras no críticas.

---

## Pendientes para Próxima Versión

### Prioridad Alta 🟡

#### 1. Tests en el frontend
- No hay tests unitarios en Svelte
- Priorizar: Podium, Celebration, AutoForm, CategoriaForm
- Herramienta propuesta: `vitest` + `@testing-library/svelte`

#### 2. Manejo de errores en apiFetch
- Actualmente no distingue entre error de red y error del backend
- Agregar timeout y retry en caso de conexión fallida

#### 3. Feedback visual en operaciones lentas
- Ya existe `actionLoading` pero no en todas las operaciones
- Unificar con un estado global de carga

### Prioridad Media 🟢

#### 4. Migraciones versionadas
- Actualmente todo en `database/sqlite.go` (archivo único)
- Implementar sistema de migraciones con números de versión
- Permitir rollback de migraciones

#### 5. Refactor del handler de fixture
- `fixture_handler.go` mezcla lógica de competencias y legacy
- Separar en: `competencia_handler.go` + `fixture_handler.go` + `archivo_handler.go`

#### 6. Eliminar redundancia en DB
- `fixtures.categoria_id` es redundante con `competencias.categoria_id`
- Requiere migración y actualización de queries

#### 7. Unificar tipado en frontend
- Muchos estados usan `any` en lugar de interfaces tipadas
- Crear tipos compartidos en `$lib/types.ts`

### Prioridad Baja 🔵

#### 8. PWA avanzado
- El service worker ya está configurado (vite-plugin-pwa)
- Mejorar estrategia de caché para assets estáticos

#### 9. Notificaciones push
- Opcional: notificar cuando una carrera finaliza
- Requiere Service Worker + backend endpoint

#### 10. Modo oscuro automático
- Detectar `prefers-color-scheme` del sistema
- Actualmente solo usa el tema guardado en localStorage

---

## Bugs Conocidos

| Bug | Estado | Notas |
|-----|--------|-------|
| Manager no persiste en segundo plano | 🟡 Abierto | El proceso se cierra al salir del shell. Usar `nohup` o `systemd` |
| API de manager ejecuta comandos como root para firewall | 🟢 Menor | Manejar con sudo o documentar |

---

## Deuda Técnica

| Item | Impacto | 
|------|:-------:|
| Handler de fixture mezcla responsabilidades | Alto |
| Sin migraciones versionadas | Medio |
| `fixtures.categoria_id` redundante | Medio |
| Frontend sin tipado fuerte | Medio |
| Sin tests de frontend | Medio |
| apiFetch sin manejo de errores de red | Bajo |

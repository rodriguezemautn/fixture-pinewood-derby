# Tareas Pendientes — Próxima Sesión

## Prioridad Crítica 🔴

### 1. Verificar flujo completo de competencias
- **Problema**: Usuario reporta que "+Nueva Competencia" aún no funciona después del fix
- **Posible causa**: La DB existente no tiene `competencia_id` en fixtures viejos
- **Solución inmediata**: `rm -f fixture.db*` y arrancar fresco
- **Solución definitiva**: Si se requieren migrar datos existentes, agregar migración SQL para fixtures que no tienen competencia:
  ```sql
  -- Crear competencias para fixtures huerfanos
  INSERT INTO competencias (id, categoria_id, numero, nombre, estado, rondas)
  SELECT f.id, f.categoria_id, 1, c.nombre || ' - Competencia #1', 'finalizada', f.rondas
  FROM fixtures f JOIN categorias c ON f.categoria_id = c.id
  WHERE f.competencia_id IS NULL;
  
  -- Vincular fixtures huerfanos
  UPDATE fixtures SET competencia_id = id WHERE competencia_id IS NULL;
  ```
- **Archivos involucrados**: `fixture_service.go`, `fixture.go` (repo), `+page.svelte` (fixture)

### 2. Validar endpoint GET /api/competencias/{id}/fixture
- **Problema**: Verificar que devuelve el fixture correcto después del fix
- **Prueba**: `curl localhost:8080/api/competencias/{id}/fixture`
- **Archivos**: `fixture.go` (repo SQLite, método GetByCompetencia)

## Prioridad Alta 🟡

### 3. Mejorar feedback visual en creación de competencia
- **Problema**: El usuario no ve confirmación de que la competencia se creó
- **Solución parcial**: Ya se agregó mensaje de éxito con timeout (commit 4091ec6)
- **Pendiente**: Verificar que se muestre correctamente

### 4. Sincronizar fixture page cuando no hay competencia seleccionada
- **Problema**: Si se navega a `/admin/categorias/{id}/fixture` sin `?competencia=xxx`, no se muestra nada
- **Solución**: Redirigir a la primera competencia disponible, o mostrar lista

### 5. Desempate selector
- **Problema**: El selector de desempate muestra checkbox por posición. Verificar que funcione correctamente con la UI actual
- **Archivos**: `+page.svelte` (fixture)

## Prioridad Media 🟢

### 6. Refresh automático en vista pública
- La vista de carreras públicas tiene polling cada 10s. Verificar que siga funcionando.

### 7. Unificar sistema de botones
- Algunos botones usan gradiente, otros color sólido. Decidir estándar.

### 8. Skeleton loaders
- Ya implementados en carrera view y landing. Verificar en otras páginas.

## Bugs Conocidos

| Bug | Estado | Notas |
|-----|--------|-------|
| fixture.competencia_id no se setea | ✅ Fixeado (aa98c54) | Requiere DB fresca |
| GenerarFinal con categoriaId | ✅ Fixeado (aa98c54) | Ahora usa competenciaId |
| Proxy Vite 8 bloquea POST | ✅ Fixeado (96ec906) | Sin changeOrigin |

## Deuda Técnica

- [ ] El campo `fixtures.categoria_id` es redundante ahora que existe `competencias.categoria_id`
- [ ] La migración ALTER TABLE para `peso` y `registrado_at` podría fallar en DBs nuevas
- [ ] No hay migraciones versionadas (todo en un mismo archivo)
- [ ] El handler de fixture mezcla lógica de competencias y legacy

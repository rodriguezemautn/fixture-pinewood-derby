# Proposal: Auth Admin

## Intent
Proteger las rutas de administración con autenticación simple (PIN) para cumplir RB1 y RB2 del SRS: admin gestiona, visualizador solo ve.

## Scope
**In**: Login endpoint con PIN, middleware de autenticación para rutas /api/admin, frontend login page, protección de rutas admin, redirección a login.
**Out**: Registro de usuarios, OAuth, permisos granular, multi-tenant.

## Capabilities
- New: `auth-admin` — Autenticación + protección de rutas admin

## Approach
PIN configurable via env (`ADMIN_PIN`). Login devuelve token JWT simple. Chi middleware verifica token en rutas /api/admin. Frontend: login page, localStorage para token, guard en layout admin.

## Files
- `backend/internal/auth/middleware.go` — Middleware JWT
- `backend/internal/handler/auth.go` — Login handler
- `backend/internal/router/router.go` — Agrupar rutas admin
- `frontend/src/routes/login/+page.svelte` — Login page
- `frontend/src/routes/admin/+layout.svelte` — Auth guard

## Success Criteria
- [ ] POST /api/auth/login con PIN correcto → 200 + token
- [ ] POST /api/auth/login con PIN incorrecto → 401
- [ ] Rutas /api/admin/* requieren token
- [ ] Rutas públicas (/api/categorias, /api/carreras) sin auth
- [ ] Frontend: login page + redirect si no autenticado
- [ ] Admin routes protegidas, públicas accesibles sin login

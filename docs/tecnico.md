# Documento Técnico — Fixture D15

> **Versión**: 1.0  
> **Última actualización**: 2026-06-25

---

## 1. Instalación

### Requisitos

```bash
go version  # ≥ 1.22
node --version  # ≥ 20
npm --version   # ≥ 10
```

### Clonar e instalar dependencias

```bash
git clone <repo-url>
cd fixture

# Instalar dependencias
make init
# o manual:
cd backend && go mod tidy
cd frontend && npm install
```

---

## 2. Desarrollo

### Iniciar entorno de desarrollo

```bash
# Todo junto (backend + frontend)
make dev

# O por separado:
make dev-backend    # Go con air (hot reload)
make dev-frontend   # Vite con HMR
```

El backend corre en `http://localhost:8080`, el frontend en `http://localhost:5173`.

### Estructura de commits

Usamos **conventional commits**:

```
feat: nueva funcionalidad
fix: corrección de bug
docs: documentación
refactor: refactorización
style: cambios de estilo/UI
test: tests
```

### Proxy Vite

El frontend en desarrollo usa Vite proxy para redirigir `/api/*` al backend:

```ts
// vite.config.ts
server: {
    host: true,
    proxy: {
        '/api': 'http://localhost:8080',
        '/uploads': 'http://localhost:8080'
    }
}
```

---

## 3. Build y Deploy

### Compilación

```bash
# Backend (binario estático)
cd backend && CGO_ENABLED=0 go build -o bin/api ./cmd/api

# Frontend (SPA estática)
cd frontend && npm run build
# → frontend/build/

# Manager (consola de gestión)
go build -o bin/manager ./manager/
```

El binario del backend es **100% estático** (CGO_ENABLED=0) y no requiere
dependencias del sistema — ni siquiera SQLite ya que usa `modernc.org/sqlite`.

### Deploy en red local

```bash
# 1. Compilar backend y frontend
make build

# 2. Iniciar backend
cd backend && ./bin/api &

# 3. Servir frontend (usar un servidor estático o abrir index.html)
cd frontend/build && python3 -m http.server 5173

# 4. Abrir puertos en firewall si es necesario
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --permanent --add-port=5173/tcp
sudo firewall-cmd --reload
```

### Acceso desde la red

```
Frontend: http://<IP>:5173
Backend:  http://<IP>:8080/api/...
Manager:  http://<IP>:9099
```

---

## 4. Testing

### Backend

```bash
cd backend && go test -v -cover ./...
```

| Paquete | Tests | Cobertura |
|---------|:-----:|:---------:|
| auth | 1 | Middleware HMAC |
| handler | 15 | Endpoints REST |
| repository/sqlite | 9 | Persistencia |
| service | 23 | Lógica de negocio + scheduler |

### Frontend

El frontend no tiene tests unitarios actualmente. Las validaciones
se realizan mediante pruebas funcionales manuales.

---

## 5. Consola de Gestión

### Script por terminal

```bash
./bin/manage.sh [comando]

Comandos:
  status    Ver estado de backend y frontend
  start     Iniciar servicios
  stop      Detener servicios
  logs      Ver logs del backend
  build     Compilar backend
  db-reset  Resetear base de datos
  ports     Abrir puertos firewall
  menu      Menú interactivo
```

### Dashboard web (independiente)

```bash
# Iniciar (puerto 9099)
./bin/manager

# O compilar desde código
cd manager && go build -o ../bin/manager .
```

El dashboard web permite:
- Monitorear estado de backend, frontend y DB
- Iniciar/detener/ reiniciar servicios
- Ver logs en vivo (backend y frontend en paralelo)
- Compilar backend
- Resetear base de datos
- Abrir puertos firewall
- Ver información del sistema (CPU, RAM, disco)

**Acceso**: `http://localhost:9099`

---

## 6. Base de Datos

### Archivo

La DB SQLite se almacena en `backend/fixture.db`.
En modo WAL (Write-Ahead Logging) genera además `fixture.db-wal` y `fixture.db-shm`.

### Consultas útiles

```bash
# Ver todas las categorías
sqlite3 backend/fixture.db "SELECT * FROM categorias;"

# Ver autos de una categoría
sqlite3 backend/fixture.db "SELECT numero, nombre, creador FROM autos WHERE categoria_id = '<ID>';"

# Ver competencias con su estado
sqlite3 backend/fixture.db "SELECT numero, nombre, estado FROM competencias;"

# Ver carreras archivadas
sqlite3 backend/fixture.db "SELECT categoria_nombre, winner_nombre, fecha FROM archivos_carrera;"
```

### Reset

```bash
rm backend/fixture.db backend/fixture.db-wal backend/fixture.db-shm
# Al reiniciar el backend, las migraciones recrean la DB vacía.
```

---

## 7. Problemas Conocidos

### Token expira después de 24h
- **Síntoma**: Error "token inválido" al hacer POST
- **Solución**: Loguearse de nuevo en `/login`
- **Comportamiento**: `apiFetch()` redirige automáticamente al login

### Puerto 9099 ocupado por Docker
- **Síntoma**: El manager no puede iniciar
- **Solución**: Cambiar el puerto en `manager/main.go` (constante `ManagerPort`)
- Alternativa: detener el contenedor que ocupa el puerto

### Proyector con bajo contraste
- **Síntoma**: Los colores se ven lavados
- **Solución**: Usar **tema claro** (botón ☀️/🌙 en el header)
- El tema claro tiene fondos más claros (`#FFF5EB`) y texto oscuro

### Las rutas admin no verifican token al cargar
- Las GET requests no requieren auth, así que las páginas admin se cargan
- El token se verifica al primer POST — si expiró, redirige a login

---

## 8. Deuda Técnica

- [ ] `fixtures.categoria_id` es redundante (ya está en `competencias.categoria_id`)
- [ ] No hay migraciones versionadas (todo en un mismo archivo SQL)
- [ ] Handler de fixture mezcla lógica de competencias y legacy
- [ ] Frontend sin tests unitarios
- [ ] Tipado laxo en states del frontend (muchos `any`)
- [ ] Las rutas admin no verifican el token al cargar la página
- [ ] El archivo `apiFetch.ts` no maneja errores de red genéricos

---

## 9. Troubleshooting

### Error: "address already in use"
```bash
# Encontrar proceso en el puerto
lsof -i :8080
# Matarlo
kill -9 <PID>
```

### Error: "no se puede eliminar una categoría con autos"
El sistema protege la integridad referencial. Primero eliminá los autos
de la categoría, después la categoría.

### Error: "no se puede eliminar un auto que fue campeón"
El auto aparece como ganador en una carrera archivada. No se puede
eliminar para preservar el historial.

### El frontend no conecta con el backend
Verificar que:
1. El backend esté corriendo en `localhost:8080`
2. El proxy de Vite esté configurado (`vite.config.ts`)
3. No haya un firewall bloqueando

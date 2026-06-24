# Delta for Infra-base

## ADDED Requirements

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

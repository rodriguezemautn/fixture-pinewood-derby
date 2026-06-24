# Gestion-categorias Specification

## Purpose

CRUD completo de categorías de carrera, con persistencia SQLite, validaciones de negocio, y endpoints REST.

**Req SRS**: RF1 — "El sistema deberá poder registrar las categorías de la fecha."

## Requirements

### Requirement: Crear categoría

El sistema **MUST** permitir crear una categoría con nombre, edad mínima y edad máxima.

#### Scenario: Creación exitosa

- GIVEN datos válidos: nombre "Pre-Juveniles", edadMin 10, edadMax 12
- WHEN se envía `POST /api/categorias` con esos datos
- THEN retorna `201 Created` con la categoría creada
- AND la categoría tiene un ID único asignado

#### Scenario: Nombre vacío rechazado

- GIVEN datos con nombre vacío
- WHEN se envía `POST /api/categorias`
- THEN retorna `400 Bad Request` con mensaje de error

#### Scenario: EdadMin mayor que EdadMax rechazado

- GIVEN datos con edadMin 15 y edadMax 10
- WHEN se envía `POST /api/categorias`
- THEN retorna `400 Bad Request`

### Requirement: Listar categorías

El sistema **MUST** retornar todas las categorías registradas.

#### Scenario: Listado exitoso

- GIVEN 3 categorías registradas
- WHEN se envía `GET /api/categorias`
- THEN retorna `200 OK` con un array de 3 categorías

### Requirement: Obtener categoría por ID

El sistema **MUST** retornar una categoría por su ID.

#### Scenario: Categoría existe

- GIVEN una categoría con ID "abc-123"
- WHEN se envía `GET /api/categorias/abc-123`
- THEN retorna `200 OK` con la categoría

#### Scenario: Categoría no existe

- GIVEN ningún ID coincide
- WHEN se envía `GET /api/categorias/inexistente`
- THEN retorna `404 Not Found`

### Requirement: Actualizar categoría

El sistema **MUST** permitir actualizar nombre, edadMin y edadMax de una categoría existente.

#### Scenario: Actualización exitosa

- GIVEN una categoría existente
- WHEN se envía `PUT /api/categorias/{id}` con nuevos datos
- THEN retorna `200 OK` con la categoría actualizada

### Requirement: Eliminar categoría

El sistema **MUST** permitir eliminar una categoría por ID.

#### Scenario: Eliminación exitosa

- GIVEN una categoría existente sin autos asociados
- WHEN se envía `DELETE /api/categorias/{id}`
- THEN retorna `204 No Content`
- AND la categoría ya no aparece en el listado

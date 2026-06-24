# Gestion-autos Specification

## Purpose

CRUD de autos de madera asociados a categorías, con foto, número de auto, nombre de creador y edad.

**Req SRS**: RF2

## Requirements

### Requirement: Crear auto en categoría

El sistema **MUST** permitir registrar un auto dentro de una categoría existente.

#### Scenario: Creación exitosa

- GIVEN una categoría con ID "cat-1"
- WHEN se envía `POST /api/categorias/cat-1/autos` con número, nombre, creador, edad y foto
- THEN retorna `201 Created` con el auto creado
- AND el auto tiene un ID único y referencia a cat-1

#### Scenario: Categoría no existe

- GIVEN ninguna categoría con ID "no-existe"
- WHEN se envía `POST /api/categorias/no-existe/autos`
- THEN retorna `404 Not Found`

#### Scenario: Número duplicado en misma categoría

- GIVEN un auto con número 5 ya registrado en cat-1
- WHEN se envía `POST /api/categorias/cat-1/autos` con número 5
- THEN retorna `400 Bad Request`

### Requirement: Listar autos por categoría

El sistema **MUST** retornar todos los autos de una categoría.

#### Scenario: Listado exitoso

- GIVEN 3 autos en cat-1
- WHEN se envía `GET /api/categorias/cat-1/autos`
- THEN retorna `200 OK` con un array de 3 autos

#### Scenario: Categoría sin autos

- GIVEN cat-2 sin autos registrados
- WHEN se envía `GET /api/categorias/cat-2/autos`
- THEN retorna `200 OK` con un array vacío

### Requirement: Obtener auto por ID

El sistema **MUST** retornar un auto individual.

#### Scenario: Auto existe

- GIVEN un auto con ID "auto-1"
- WHEN se envía `GET /api/autos/auto-1`
- THEN retorna `200 OK` con el auto

#### Scenario: Auto no existe

- WHEN se envía `GET /api/autos/inexistente`
- THEN retorna `404 Not Found`

### Requirement: Actualizar auto

El sistema **MUST** permitir actualizar los datos de un auto.

#### Scenario: Actualización exitosa

- GIVEN un auto existente
- WHEN se envía `PUT /api/autos/{id}` con nuevos datos
- THEN retorna `200 OK` con el auto actualizado

### Requirement: Eliminar auto

El sistema **MUST** permitir eliminar un auto.

#### Scenario: Eliminación exitosa

- GIVEN un auto existente
- WHEN se envía `DELETE /api/autos/{id}`
- THEN retorna `204 No Content`

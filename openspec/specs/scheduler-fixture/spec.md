# Scheduler-fixture Specification

## Requirements

### Requirement: Generar fixture

El sistema **MUST** generar un fixture de R rondas dado un conjunto de autos en una categoría.

#### Scenario: Fixture generado
- GIVEN 8 autos en categoría "Pre-Juveniles"
- WHEN se solicita generar fixture con 3 rondas
- THEN se crean 6 heats de clasificación (2 por ronda)
- AND cada auto aparece exactamente 3 veces en los heats
- AND los heats tienen máximo 4 autos cada uno

### Requirement: Registrar resultado

El sistema **MUST** registrar el orden de llegada de cada heat.

#### Scenario: Resultado registrado
- GIVEN un heat con 4 autos
- WHEN se registra orden: [auto-3, auto-1, auto-4, auto-2]
- THEN auto-3 recibe 5 puntos, auto-1 recibe 3, auto-4 recibe 2, auto-2 recibe 1

### Requirement: Calcular posiciones

El sistema **MUST** calcular y retornar la tabla de posiciones acumuladas.

#### Scenario: Posiciones calculadas
- GIVEN 2 heats completados en una categoría
- WHEN se consulta la tabla de posiciones
- THEN retorna autos ordenados por puntos descendente

### Requirement: Seleccionar top 4 para final

El sistema **MUST** seleccionar los 4 autos con mayor puntaje para la carrera final.

#### Scenario: Final generada
- GIVEN todos los heats de clasificación completados
- WHEN se solicita generar final
- THEN se crea una carrera final con los 4 autos de mayor puntaje
- AND se rompen empates por: 1) mayor cantidad de 1° puestos, 2) menor edad

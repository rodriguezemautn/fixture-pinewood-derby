# Exploration: Scheduler de Fixture

## El Problema

Dado N autos en una categoría, generar una secuencia de carreras de hasta 4 autos cada una que:
1. Determine los 4 mejores autos para la carrera final
2. Sea **óptima** (mínimas carreras necesarias)
3. Sea **equitativa** (misma cantidad de carreras por auto, variedad de oponentes)

## Algoritmo Propuesto: Round-Robin por Grupos

### Fase 1: Heat de Clasificación

Cada auto corre R rondas. En cada ronda se dividen en heats de 4 autos (el último heat puede tener menos).

**Puntaje por heat**:
| Posición | Puntos |
|----------|--------|
| 1° | 5 |
| 2° | 3 |
| 3° | 2 |
| 4° | 1 |
| DNS | 0 |

### Fase 2: Final

Los 4 autos con mayor puntaje acumulado corren la carrera final.

### Cantidad de Rondas

Para N autos, con R rondas:
- Heats por ronda: `ceil(N / 4)`
- Total heats clasificación: `R * ceil(N / 4)`
- Cada auto corre: `R` carreras

R recomendado: suficiente para que cada auto corra al menos 3 veces.

| N autos | Rondas | Heats clasif | Total carreras |
|---------|--------|--------------|----------------|
| 8 | 3 | 6 | 7 |
| 12 | 3 | 9 | 10 |
| 16 | 3 | 12 | 13 |
| 20 | 4 | 20 | 21 |

### Asignación de Heats (Equidad)

Para maximizar equidad:
1. **Round 1**: Orden aleatorio
2. **Round 2+**: Swiss-system ligero — agrupar por performance (puntaje acumulado) para emparejar autos de nivel similar

Esto evita que un auto domine heats débiles, dando mediciones más precisas.

## Data Model

```go
type Fixture struct {
    ID           string
    CategoriaID  string
    Autos        []string // IDs de autos participantes
    Carreras     []Carrera // heats generados
    Final        *Carrera  // carrera final
    Estado       string   // "pendiente", "en_curso", "finalizado"
}

type Resultado struct {
    AutoID string
    Posicion int
    Puntos  int
    CarreraID string
}
```

## Approach Técnico

1. Service `FixtureService` con método `GenerarFixture(categoriaID string, rondas int)` 
2. Store resultados de cada heat
3. Calculate standings after each heat
4. Algoritmo puro en Go (sin dependencias)
5. API endpoints para gestionar fixture y registrar resultados

# Design: Scheduler de Fixture

## Algorithm

```
GenerarFixture(autos, rondas):
  1. Shuffle autos aleatoriamente
  2. Para cada ronda 1..R:
     a. Ordenar autos por puntaje acumulado (ascendente)
     b. Dividir en grupos de 4
     c. Asignar cada grupo a un heat
  3. Retornar lista de heats generados

CalcularPuntos(posicion):
  1°=5, 2°=3, 3°=2, 4°=1, DNS=0

SeleccionarFinal(standings):
  Top 4 por puntos
  Desempate: más 1° puestos → menor edad
```

## Files
- `backend/internal/domain/fixture.go` — types
- `backend/internal/service/scheduler.go` — algoritmo puro
- `backend/internal/repository/sqlite/fixture.go` — repo SQLite
- `backend/internal/service/fixture_service.go` — service orquestador
- `backend/internal/handler/fixture.go` — API REST
- `frontend/src/routes/admin/categorias/[id]/fixture/` — Frontend

## API
- `POST /api/categorias/{id}/fixture?rondas=3` → generar fixture
- `GET /api/categorias/{id}/fixture` → ver fixture
- `GET /api/categorias/{id}/posiciones` → tabla de posiciones
- `POST /api/carreras/{id}/resultado` → registrar orden llegada
- `POST /api/categorias/{id}/final` → generar carrera final

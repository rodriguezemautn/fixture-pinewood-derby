// Package service contiene el algoritmo de scheduler de fixture.
package service

import (
	"math/rand"
	"sort"

	"github.com/ema/fixture/backend/internal/domain"
)

// PuntosPorPosicion asigna puntos según la posición en el heat.
var PuntosPorPosicion = map[int]int{1: 5, 2: 3, 3: 2, 4: 1}

// GenerarHeats genera heats de clasificación para N autos con R rondas.
// Cada ronda MEZCLA los autos para que los grupos sean siempre distintos.
// Cuando hay resultados, usa Swiss-system (agrupa por performance similar).
func GenerarHeats(autoIDs []string, rondas int) []domain.Heat {
	if len(autoIDs) == 0 || rondas < 1 {
		return nil
	}

	n := len(autoIDs)
	puntos := make(map[string]int)
	for _, id := range autoIDs {
		puntos[id] = 0
	}

	var heats []domain.Heat
	heatNum := 0

	for ronda := 0; ronda < rondas; ronda++ {
		ordenados := make([]string, n)
		copy(ordenados, autoIDs)

		// Si hay puntos (resultados de rondas anteriores), ordenar por rendimiento
		// Sino, mezclar aleatoriamente para que cada ronda sea distinta
		tienePuntos := false
		for _, p := range puntos {
			if p > 0 {
				tienePuntos = true
				break
			}
		}

		if tienePuntos {
			// Swiss-system: agrupar autos de rendimiento similar
			sort.SliceStable(ordenados, func(i, j int) bool {
				pi := puntos[ordenados[i]]
				pj := puntos[ordenados[j]]
				if pi != pj {
					return pi < pj
				}
				// Mismos puntos: shuffle para variar
				return rand.Intn(2) == 0
			})

			// Mezcla parcial intra-grupo de performance para variar carriles
			for i := 0; i < n; i += 4 {
				end := i + 4
				if end > n {
					end = n
				}
				sub := ordenados[i:end]
				if len(sub) > 1 {
					rand.Shuffle(len(sub), func(i, j int) {
						sub[i], sub[j] = sub[j], sub[i]
					})
				}
			}
		} else {
			// Sin resultados: mezcla completa para variedad entre rondas
			rand.Shuffle(len(ordenados), func(i, j int) {
				ordenados[i], ordenados[j] = ordenados[j], ordenados[i]
			})
		}

		// Dividir en grupos de hasta 4
		for i := 0; i < n; i += 4 {
			end := i + 4
			if end > n {
				end = n
			}
			grupo := make([]string, end-i)
			copy(grupo, ordenados[i:end])

			heatNum++
			heats = append(heats, domain.Heat{
				Numero:     heatNum,
				AutoIDs:    grupo,
				Completado: false,
			})
		}
	}

	return heats
}

// CalcularStandings calcula la tabla de posiciones a partir de heats completados.
func CalcularStandings(autos map[string]*domain.Auto, heats []domain.Heat) []domain.Standing {
	standingsMap := make(map[string]*domain.Standing)
	for _, auto := range autos {
		standingsMap[auto.ID] = &domain.Standing{
			AutoID:     auto.ID,
			Nombre:     auto.Nombre,
			Numero:     auto.Numero,
			Puntos:     0,
			Posiciones: make(map[int]int),
			Carreras:   0,
		}
	}

	for _, heat := range heats {
		if !heat.Completado || len(heat.OrdenLlegada) == 0 {
			continue
		}
		for pos, autoID := range heat.OrdenLlegada {
			puesto := pos + 1
			if s, ok := standingsMap[autoID]; ok {
				s.Puntos += PuntosPorPosicion[puesto]
				s.Posiciones[puesto]++
				s.Carreras++
			}
		}
	}

	standings := make([]domain.Standing, 0, len(standingsMap))
	for _, s := range standingsMap {
		standings = append(standings, *s)
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Puntos != standings[j].Puntos {
			return standings[i].Puntos > standings[j].Puntos
		}
		if standings[i].Posiciones[1] != standings[j].Posiciones[1] {
			return standings[i].Posiciones[1] > standings[j].Posiciones[1]
		}
		a1 := autos[standings[i].AutoID]
		a2 := autos[standings[j].AutoID]
		if a1 != nil && a2 != nil && a1.Edad != a2.Edad {
			return a1.Edad < a2.Edad
		}
		return standings[i].AutoID < standings[j].AutoID
	})

	return standings
}

// SeleccionarFinal selecciona los top 4 autos para la carrera final.
func SeleccionarFinal(standings []domain.Standing) []string {
	n := 4
	if len(standings) < 4 {
		n = len(standings)
	}
	autoIDs := make([]string, n)
	for i := 0; i < n; i++ {
		autoIDs[i] = standings[i].AutoID
	}
	return autoIDs
}

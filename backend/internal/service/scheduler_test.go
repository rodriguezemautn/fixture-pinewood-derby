package service

import (
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
)

func TestGenerarHeats_8autos_3rondas(t *testing.T) {
	autos := []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8"}
	heats := GenerarHeats(autos, 3)

	// 8 autos / 4 = 2 heats por ronda * 3 rondas = 6 heats
	if len(heats) != 6 {
		t.Errorf("expected 6 heats, got %d", len(heats))
	}

	// Cada auto debe aparecer exactamente 3 veces (3 rondas)
	apariciones := make(map[string]int)
	for _, h := range heats {
		for _, id := range h.AutoIDs {
			apariciones[id]++
		}
	}
	for _, id := range autos {
		if apariciones[id] != 3 {
			t.Errorf("auto %s appeared %d times, expected 3", id, apariciones[id])
		}
	}

	// Ningún heat debe tener más de 4 autos
	for _, h := range heats {
		if len(h.AutoIDs) > 4 {
			t.Errorf("heat %d has %d autos, max 4", h.Numero, len(h.AutoIDs))
		}
	}
}

func TestGenerarHeats_12autos_3rondas(t *testing.T) {
	autos := make([]string, 12)
	for i := range autos {
		autos[i] = string(rune('a' + i))
	}

	heats := GenerarHeats(autos, 3)

	// 12/4 = 3 heats por ronda * 3 = 9 heats
	if len(heats) != 9 {
		t.Errorf("expected 9 heats, got %d", len(heats))
	}

	// Cada auto aparece 3 veces
	apariciones := make(map[string]int)
	for _, h := range heats {
		for _, id := range h.AutoIDs {
			apariciones[id]++
		}
	}
	for _, id := range autos {
		if apariciones[id] != 3 {
			t.Errorf("auto %s appeared %d times", id, apariciones[id])
		}
	}
}

func TestGenerarHeats_3autos_2rondas(t *testing.T) {
	autos := []string{"a1", "a2", "a3"}
	heats := GenerarHeats(autos, 2)

	// 3 autos = 1 heat por ronda (con 3 autos)
	if len(heats) != 2 {
		t.Errorf("expected 2 heats, got %d", len(heats))
	}

	for _, h := range heats {
		if len(h.AutoIDs) != 3 {
			t.Errorf("heat %d: expected 3 autos, got %d", h.Numero, len(h.AutoIDs))
		}
	}
}

func TestGenerarHeats_distintasRondas(t *testing.T) {
	autos := []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8"}
	heats := GenerarHeats(autos, 3)

	// 8 autos = 2 heats por ronda, 3 rondas = 6 heats
	if len(heats) != 6 {
		t.Fatalf("expected 6 heats, got %d", len(heats))
	}

	// Verificar que las rondas NO tengan los mismos grupos
	// Ronda 1: heats 0-1, Ronda 2: heats 2-3, Ronda 3: heats 4-5
	ronda1 := [][]string{heats[0].AutoIDs, heats[1].AutoIDs}
	ronda2 := [][]string{heats[2].AutoIDs, heats[3].AutoIDs}

	// Aplanar ambas rondas y verificar que el orden sea distinto
	orden1 := make([]string, 0, 8)
	for _, h := range ronda1 {
		orden1 = append(orden1, h...)
	}
	orden2 := make([]string, 0, 8)
	for _, h := range ronda2 {
		orden2 = append(orden2, h...)
	}

	// Verificar que las rondas no sean identicas (deberian ser diferentes por el shuffle)
	mismaRonda := true
	for i := 0; i < 8; i++ {
		if orden1[i] != orden2[i] {
			mismaRonda = false
			break
		}
	}
	if mismaRonda {
		t.Error("all rounds produced identical heat compositions — shuffle not working")
	}
}

func TestGenerarHeats_empty(t *testing.T) {
	if h := GenerarHeats(nil, 3); h != nil {
		t.Error("expected nil for nil input")
	}
	if h := GenerarHeats([]string{}, 3); h != nil {
		t.Error("expected nil for empty input")
	}
	if h := GenerarHeats([]string{"a1"}, 0); h != nil {
		t.Error("expected nil for 0 rounds")
	}
}

func TestCalcularStandings_basic(t *testing.T) {
	autos := map[string]*domain.Auto{
		"a1": {ID: "a1", Nombre: "Turbo", Numero: 1, Edad: 10},
		"a2": {ID: "a2", Nombre: "Rayo", Numero: 2, Edad: 11},
		"a3": {ID: "a3", Nombre: "Flash", Numero: 3, Edad: 12},
		"a4": {ID: "a4", Nombre: "Rocket", Numero: 4, Edad: 13},
	}

	heats := []domain.Heat{
		{Completado: true, OrdenLlegada: []string{"a3", "a1", "a2", "a4"}},
		{Completado: true, OrdenLlegada: []string{"a1", "a4", "a2", "a3"}},
	}

	standings := CalcularStandings(autos, heats)

	if len(standings) != 4 {
		t.Fatalf("expected 4 standings, got %d", len(standings))
	}

	// a1: 3+5=8, a3: 5+1=6, a2: 2+2=4, a4: 1+3=4
	// a1 first, a3 second, a2/a4 tied → a2 wins by AutoID
	if standings[0].AutoID != "a1" {
		t.Errorf("expected a1 first (8pts), got %s", standings[0].AutoID)
	}
	if standings[0].Puntos != 8 || standings[1].Puntos != 6 {
		t.Errorf("expected points [8,6,4,4], got [%d,%d,%d,%d]",
			standings[0].Puntos, standings[1].Puntos, standings[2].Puntos, standings[3].Puntos)
	}
}

func TestCalcularStandings_tiebreak(t *testing.T) {
	autos := map[string]*domain.Auto{
		"a1": {ID: "a1", Nombre: "A", Numero: 1, Edad: 10},
		"a2": {ID: "a2", Nombre: "B", Numero: 2, Edad: 15},
	}

	heats := []domain.Heat{
		{Completado: true, OrdenLlegada: []string{"a1", "a2"}},
		{Completado: true, OrdenLlegada: []string{"a2", "a1"}},
	}

	standings := CalcularStandings(autos, heats)

	// Both have 5+3=8 points, a1 has 1 first place, a2 has 1 first place
	// Tiebreak: a1 is younger (10 vs 15) → a1 wins
	if standings[0].AutoID != "a1" {
		t.Errorf("expected a1 first (younger), got %s", standings[0].AutoID)
	}
}

func TestCalcularStandings_ignoresUncompleted(t *testing.T) {
	autos := map[string]*domain.Auto{
		"a1": {ID: "a1", Nombre: "A", Numero: 1, Edad: 10},
		"a2": {ID: "a2", Nombre: "B", Numero: 2, Edad: 11},
	}

	heats := []domain.Heat{
		{Completado: false, OrdenLlegada: []string{"a1", "a2"}},          // should be ignored
		{Completado: true, OrdenLlegada: []string{"a2", "a1"}},           // counted
		{Completado: true, OrdenLlegada: []string{}},                     // empty, should be ignored
	}

	standings := CalcularStandings(autos, heats)
	if standings[0].Puntos != 5 {
		t.Errorf("expected first place to have 5 points (only 1 heat counted), got %d", standings[0].Puntos)
	}
}

func TestSeleccionarFinal_top4(t *testing.T) {
	standings := []domain.Standing{
		{AutoID: "a1", Puntos: 20},
		{AutoID: "a2", Puntos: 18},
		{AutoID: "a3", Puntos: 15},
		{AutoID: "a4", Puntos: 12},
		{AutoID: "a5", Puntos: 10},
	}

	final := SeleccionarFinal(standings)
	if len(final) != 4 {
		t.Errorf("expected 4 finalists, got %d", len(final))
	}
	if final[0] != "a1" || final[3] != "a4" {
		t.Errorf("expected [a1,a2,a3,a4], got %v", final)
	}
}

func TestSeleccionarFinal_lessThan4(t *testing.T) {
	standings := []domain.Standing{
		{AutoID: "a1", Puntos: 5},
		{AutoID: "a2", Puntos: 3},
	}

	final := SeleccionarFinal(standings)
	if len(final) != 2 {
		t.Errorf("expected 2 finalists, got %d", len(final))
	}
}

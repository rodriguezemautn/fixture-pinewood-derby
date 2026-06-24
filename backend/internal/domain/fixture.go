package domain

// Fixture representa un fixture de carreras para una categoría.
type Fixture struct {
	ID          string
	CategoriaID string
	Estado      string   // "pendiente", "en_curso", "finalizado"
	Rondas      int
	Heats       []Heat
	Final       *Carrera
}

// Heat representa una carrera individual en el fixture.
type Heat struct {
	ID       string
	Numero   int
	AutoIDs  []string // IDs de los autos en esta heat
	Completado bool
	OrdenLlegada []string
}

// Standing representa la posición de un auto en la tabla general.
type Standing struct {
	AutoID    string
	Nombre    string
	Numero    int
	Puntos    int
	Posiciones map[int]int // posición -> cantidad (ej: {1: 2, 2: 1, 3: 0, 4: 1})
	Carreras  int
}

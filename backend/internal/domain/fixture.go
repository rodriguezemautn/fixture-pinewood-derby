package domain

// Fixture representa un fixture de carreras para una categoría.
type Fixture struct {
	ID            string `json:"id"`
	CategoriaID   string `json:"categoria_id"`
	CompetenciaID string `json:"competencia_id,omitempty"`
	Estado        string `json:"estado"`
	Rondas        int    `json:"rondas"`
	Heats         []Heat `json:"heats"`
	Final         *Carrera `json:"final,omitempty"`
}

// Heat representa una carrera individual en el fixture.
type Heat struct {
	ID           string   `json:"id"`
	Numero       int      `json:"numero"`
	AutoIDs      []string `json:"auto_ids"`
	Completado   bool     `json:"completado"`
	OrdenLlegada []string `json:"orden_llegada,omitempty"`
}

// Standing representa la posición de un auto en la tabla general.
type Standing struct {
	AutoID    string   `json:"auto_id"`
	Nombre    string   `json:"nombre"`
	Numero    int      `json:"numero"`
	Puntos    int      `json:"puntos"`
	Posiciones map[int]int `json:"posiciones"`
	Carreras  int      `json:"carreras"`
}

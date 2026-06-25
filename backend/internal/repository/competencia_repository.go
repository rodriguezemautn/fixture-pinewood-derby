package repository

// Competencia define una serie de carreras dentro de una categoría.
type Competencia struct {
	ID          string `json:"id"`
	CategoriaID string `json:"categoria_id"`
	Numero      int    `json:"numero"`
	Nombre      string `json:"nombre"`
	Estado      string `json:"estado"`
	Rondas      int    `json:"rondas"`
	CreatedAt   string `json:"created_at"`
}

// CompetenciaRepository define operaciones sobre competencias.
type CompetenciaRepository interface {
	Create(c *Competencia) error
	ListByCategoria(categoriaID string) ([]Competencia, error)
	GetByID(id string) (*Competencia, error)
	SetEstado(id, estado string) error
	SetNombre(id, nombre string) error
}

// Package domain contiene las entidades del negocio, libres de dependencias externas.
package domain

// Categoria representa una categoría de carrera por edad.
type Categoria struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	EdadMin   int    `json:"edad_min"`
	EdadMax   int    `json:"edad_max"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Auto representa un auto de madera registrado en una categoría.
type Auto struct {
	ID          string `json:"id"`
	CategoriaID string `json:"categoria_id"`
	Numero      int    `json:"numero"`
	Nombre      string `json:"nombre"`
	Creador     string `json:"creador"`
	Edad        int    `json:"edad"`
	Peso        int    `json:"peso"`
	FotoURL     string `json:"foto_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Carrera representa una instancia de carrera con hasta 4 autos.
type Carrera struct {
	ID           string   `json:"id"`
	CategoriaID  string   `json:"categoria_id"`
	AutoIDs      []string `json:"auto_ids"`
	OrdenLlegada []string `json:"orden_llegada"`
}

// Package domain contiene las entidades del negocio, libres de dependencias externas.
package domain

// Categoria representa una categoría de carrera por edad.
type Categoria struct {
	ID        string
	Nombre    string
	EdadMin   int
	EdadMax   int
	CreatedAt string
	UpdatedAt string
}

// Auto representa un auto de madera registrado en una categoría.
type Auto struct {
	ID          string
	Numero      int
	Nombre      string
	Creador     string
	Edad        int
	FotoURL     string
	CategoriaID string
}

// Carrera representa una instancia de carrera con hasta 4 autos.
type Carrera struct {
	ID           string
	CategoriaID  string
	AutoIDs      []string // IDs de los autos participantes (máx 4)
	OrdenLlegada []string // IDs de autos en orden de llegada
}

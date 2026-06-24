package service

import (
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
)

// autoMock implementa AutoRepository para tests.
type autoMock struct {
	autos map[string]*domain.Auto
}

func (m *autoMock) ListAll() ([]domain.Auto, error) {
	var res []domain.Auto
	for _, a := range m.autos {
		res = append(res, *a)
	}
	return res, nil
}

func (m *autoMock) ListByCategoria(categoriaID string) ([]domain.Auto, error) {
	var res []domain.Auto
	for _, a := range m.autos {
		if a.CategoriaID == categoriaID {
			res = append(res, *a)
		}
	}
	return res, nil
}
func (m *autoMock) GetByID(id string) (*domain.Auto, error) {
	a, ok := m.autos[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}
func (m *autoMock) Save(a *domain.Auto) error {
	m.autos[a.ID] = a
	return nil
}
func (m *autoMock) Update(a *domain.Auto) error {
	m.autos[a.ID] = a
	return nil
}
func (m *autoMock) Delete(id string) error {
	delete(m.autos, id)
	return nil
}
func (m *autoMock) ExistsByNumero(categoriaID string, numero int) (bool, error) {
	for _, a := range m.autos {
		if a.CategoriaID == categoriaID && a.Numero == numero {
			return true, nil
		}
	}
	return false, nil
}

// categoriaMock implementa CategoriaRepository para tests.
type categoriaMock struct {
	data map[string]*domain.Categoria
}

func (m *categoriaMock) List() ([]domain.Categoria, error) { return nil, nil }
func (m *categoriaMock) GetByID(id string) (*domain.Categoria, error) {
	c, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}
func (m *categoriaMock) Save(c *domain.Categoria) error  { return nil }
func (m *categoriaMock) Update(c *domain.Categoria) error { return nil }
func (m *categoriaMock) Delete(id string) error           { return nil }

func newAutoMocks() (*autoMock, *categoriaMock) {
	return &autoMock{autos: make(map[string]*domain.Auto)},
		&categoriaMock{data: map[string]*domain.Categoria{"cat-1": {ID: "cat-1", Nombre: "Test"}}}
}

func TestAutoService_Create_Success(t *testing.T) {
	autoR, catR := newAutoMocks()
	svc := NewAutoService(autoR, catR)

	a, err := svc.Create("cat-1", 1, "Turbo", "Juan", 10, "")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if a.Nombre != "Turbo" {
		t.Errorf("expected Turbo, got %s", a.Nombre)
	}
	if a.Numero != 1 {
		t.Errorf("expected numero 1, got %d", a.Numero)
	}
}

func TestAutoService_Create_CategoriaNotFound(t *testing.T) {
	autoR, _ := newAutoMocks()
	catR := &categoriaMock{data: make(map[string]*domain.Categoria)}
	svc := NewAutoService(autoR, catR)

	_, err := svc.Create("no-existe", 1, "Test", "J", 10, "")
	if err == nil {
		t.Error("expected error for nonexistent categoria")
	}
}

func TestAutoService_Create_DuplicateNumero(t *testing.T) {
	autoR, catR := newAutoMocks()
	svc := NewAutoService(autoR, catR)

	svc.Create("cat-1", 1, "A", "J", 10, "")
	_, err := svc.Create("cat-1", 1, "B", "J", 10, "")
	if err == nil {
		t.Error("expected error for duplicate numero")
	}
}

func TestAutoService_Create_Validation(t *testing.T) {
	autoR, catR := newAutoMocks()
	svc := NewAutoService(autoR, catR)

	tests := []struct {
		name    string
		numero  int
		nombre  string
		creador string
		edad    int
	}{
		{"nombre vacío", 1, "", "J", 10},
		{"creador vacío", 1, "A", "", 10},
		{"numero < 1", 0, "A", "J", 10},
		{"edad < 1", 1, "A", "J", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create("cat-1", tt.numero, tt.nombre, tt.creador, tt.edad, "")
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestAutoService_Update_Success(t *testing.T) {
	autoR, catR := newAutoMocks()
	svc := NewAutoService(autoR, catR)

	created, _ := svc.Create("cat-1", 1, "Old", "J", 10, "")
	updated, err := svc.Update(created.ID, 2, "New", "M", 12, "foto.jpg")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Nombre != "New" || updated.Numero != 2 {
		t.Errorf("update not applied: %+v", updated)
	}
}

func TestAutoService_Delete(t *testing.T) {
	autoR, catR := newAutoMocks()
	svc := NewAutoService(autoR, catR)

	created, _ := svc.Create("cat-1", 1, "Test", "J", 10, "")
	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	got, _ := svc.GetByID(created.ID)
	if got != nil {
		t.Error("expected nil after delete")
	}
}

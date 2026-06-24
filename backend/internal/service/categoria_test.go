package service

import (
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
)

// mockCategoriaRepo implementa repository.CategoriaRepository en memoria.
type mockCategoriaRepo struct {
	data map[string]*domain.Categoria
}

func (m *mockCategoriaRepo) List() ([]domain.Categoria, error) {
	var res []domain.Categoria
	for _, c := range m.data {
		res = append(res, *c)
	}
	return res, nil
}

func (m *mockCategoriaRepo) GetByID(id string) (*domain.Categoria, error) {
	c, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockCategoriaRepo) Save(c *domain.Categoria) error {
	m.data[c.ID] = c
	return nil
}

func (m *mockCategoriaRepo) Update(c *domain.Categoria) error {
	if _, ok := m.data[c.ID]; !ok {
		return nil
	}
	m.data[c.ID] = c
	return nil
}

func (m *mockCategoriaRepo) Delete(id string) error {
	delete(m.data, id)
	return nil
}

func newMockRepo() *mockCategoriaRepo {
	return &mockCategoriaRepo{data: make(map[string]*domain.Categoria)}
}

func TestCategoriaService_Create_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	c, err := svc.Create("Pre-Juveniles", 10, 12)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if c.Nombre != "Pre-Juveniles" {
		t.Errorf("expected Pre-Juveniles, got %s", c.Nombre)
	}
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(repo.data) != 1 {
		t.Errorf("expected 1 categoria in repo, got %d", len(repo.data))
	}
}

func TestCategoriaService_Create_Validation(t *testing.T) {
	svc := NewCategoriaService(newMockRepo())

	tests := []struct {
		name    string
		nombre  string
		edadMin int
		edadMax int
	}{
		{"nombre vacío", "", 10, 12},
		{"edad min < 1", "Test", 0, 12},
		{"edad max < 1", "Test", 10, 0},
		{"edad min >= max", "Test", 15, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(tt.nombre, tt.edadMin, tt.edadMax)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestCategoriaService_List(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	svc.Create("A", 1, 5)
	svc.Create("B", 6, 10)

	list, err := svc.List()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}

func TestCategoriaService_GetByID(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	created, _ := svc.Create("Test", 1, 5)
	got, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("getbyid failed: %v", err)
	}
	if got == nil || got.Nombre != "Test" {
		t.Errorf("expected Test, got %v", got)
	}

	notFound, _ := svc.GetByID("no-existe")
	if notFound != nil {
		t.Error("expected nil for nonexistent")
	}
}

func TestCategoriaService_Update_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	created, _ := svc.Create("Viejo", 10, 12)
	updated, err := svc.Update(created.ID, "Nuevo", 13, 15)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Nombre != "Nuevo" || updated.EdadMin != 13 {
		t.Errorf("update not applied: %+v", updated)
	}
}

func TestCategoriaService_Update_Validation(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	created, _ := svc.Create("Test", 1, 5)
	_, err := svc.Update(created.ID, "", 1, 5)
	if err == nil {
		t.Error("expected validation error on update, got nil")
	}
}

func TestCategoriaService_Delete(t *testing.T) {
	repo := newMockRepo()
	svc := NewCategoriaService(repo)

	created, _ := svc.Create("Test", 1, 5)
	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	got, _ := svc.GetByID(created.ID)
	if got != nil {
		t.Error("expected nil after delete")
	}
}

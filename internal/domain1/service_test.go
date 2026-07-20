package domain1_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/username/go-api-kit/internal/domain1"
)

type mockRepository struct {
	items []domain1.Domain1
}

func (m *mockRepository) FindAll(ctx context.Context) ([]domain1.Domain1, error) {
	return m.items, nil
}

func (m *mockRepository) FindByID(ctx context.Context, id int64) (*domain1.Domain1, error) {
	for _, item := range m.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, d *domain1.Domain1) error {
	d.ID = int64(len(m.items) + 1)
	m.items = append(m.items, *d)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, d *domain1.Domain1) error {
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestService_GetAll(t *testing.T) {
	repo := &mockRepository{
		items: []domain1.Domain1{
			{ID: 1, Name: "Item 1", CreatedAt: time.Now()},
			{ID: 2, Name: "Item 2", CreatedAt: time.Now()},
		},
	}
	svc := domain1.NewService(repo)

	items, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	repo := &mockRepository{items: []domain1.Domain1{}}
	svc := domain1.NewService(repo)

	_, err := svc.GetByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, domain1.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

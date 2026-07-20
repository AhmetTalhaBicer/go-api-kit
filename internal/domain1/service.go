package domain1

import (
	"context"
	"fmt"
	"time"
)

// Service implements the business logic layer for Domain1.
// It depends on the Repository interface and never accesses the database directly.
type Service struct {
	repo Repository
}

// NewService creates a new Domain1 Service with the provided repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetAll returns all Domain1 records.
func (s *Service) GetAll(ctx context.Context) ([]Domain1, error) {
	items, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("domain1: GetAll: %w", err)
	}
	return items, nil
}

// GetByID returns the Domain1 record with the given ID, or ErrNotFound.
func (s *Service) GetByID(ctx context.Context, id int64) (*Domain1, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("domain1: GetByID: %w", err)
	}
	if item == nil {
		return nil, ErrNotFound
	}
	return item, nil
}

// Create persists a new Domain1 record.
func (s *Service) Create(ctx context.Context, name string) (*Domain1, error) {
	now := time.Now()
	d := &Domain1{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("domain1: Create: %w", err)
	}
	return d, nil
}

// Update fetches the existing record, updates its fields, and persists the change.
func (s *Service) Update(ctx context.Context, id int64, name string) (*Domain1, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("domain1: Update: %w", err)
	}

	existing.Name = name
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("domain1: Update: %w", err)
	}
	return existing, nil
}

// Delete checks that the record exists and removes it.
func (s *Service) Delete(ctx context.Context, id int64) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return fmt.Errorf("domain1: Delete: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("domain1: Delete: %w", err)
	}
	return nil
}

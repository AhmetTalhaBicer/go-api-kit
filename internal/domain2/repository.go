package domain2

import "context"

// Repository defines the data access contract for Domain2 entities.
// Provide a concrete implementation (PostgreSQL, MySQL, in-memory, etc.)
// and inject it via NewService. This keeps the service layer database-agnostic
// and makes unit testing straightforward with a mock.
type Repository interface {
	FindAll(ctx context.Context) ([]Domain2, error)
	FindByID(ctx context.Context, id int64) (*Domain2, error)
	Create(ctx context.Context, d *Domain2) error
	Update(ctx context.Context, d *Domain2) error
	Delete(ctx context.Context, id int64) error
}

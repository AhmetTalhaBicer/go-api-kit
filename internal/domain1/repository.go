package domain1

import "context"

// Repository defines the data access contract for Domain1 entities.
// Provide a concrete implementation (PostgreSQL, MySQL, in-memory, etc.)
// and inject it via NewService. This keeps the service layer database-agnostic
// and makes unit testing straightforward with a mock.
type Repository interface {
	FindAll(ctx context.Context) ([]Domain1, error)
	FindByID(ctx context.Context, id int64) (*Domain1, error)
	Create(ctx context.Context, d *Domain1) error
	Update(ctx context.Context, d *Domain1) error
	Delete(ctx context.Context, id int64) error
}

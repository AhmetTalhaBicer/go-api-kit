package domain1

import (
	"errors"
	"time"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("domain1: not found")

// Domain1 is a placeholder entity for your first domain.
// Rename this type and its fields to match your business concept (e.g. User, Order, Post).
type Domain1 struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

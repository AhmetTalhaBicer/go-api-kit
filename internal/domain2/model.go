package domain2

import (
	"errors"
	"time"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("domain2: not found")

// Domain2 is a placeholder entity for your second domain.
// Rename this type and its fields to match your business concept (e.g. Product, Article, Category).
type Domain2 struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

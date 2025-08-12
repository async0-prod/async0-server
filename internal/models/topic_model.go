package models

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

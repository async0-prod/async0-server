package models

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Link          string    `json:"link,omitempty"`
	Author        string    `json:"author,omitempty"`
	TotalProblems int       `json:"total_problems"`
	IsActive      bool      `json:"is_active"`
	DisplayOrder  int       `json:"display_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListBasic struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

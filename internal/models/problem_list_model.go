package models

import (
	"time"

	"github.com/google/uuid"
)

type ListProblem struct {
	ListID     uuid.UUID `json:"list_id"`
	ProblemID  uuid.UUID `json:"problem_id"`
	Position   int       `json:"position"`
	IsRequired bool      `json:"is_required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

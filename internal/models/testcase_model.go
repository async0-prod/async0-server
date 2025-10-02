package models

import (
	"time"

	"github.com/google/uuid"
)

type Testcase struct {
	ID        uuid.UUID `json:"id"`
	ProblemID uuid.UUID `json:"problem_id"`
	UI        string    `json:"ui"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Position  int       `json:"position"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Slug                  string    `json:"slug"`
	Description           string    `json:"description"`
	Link                  string    `json:"link,omitempty"`
	ProblemNumber         *int      `json:"problem_number,omitempty"`
	Difficulty            string    `json:"difficulty"`
	StarterCode           any       `json:"starter_code"`
	SolutionCode          any       `json:"solution_code,omitempty"`
	TimeLimit             int       `json:"time_limit"`
	MemoryLimit           int       `json:"memory_limit"`
	AcceptanceRate        *float64  `json:"acceptance_rate,omitempty"`
	TotalSubmissions      int       `json:"total_submissions"`
	SuccessfulSubmissions int       `json:"successful_submissions"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

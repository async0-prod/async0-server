package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ProblemID       uuid.UUID `json:"problem_id"`
	Code            string    `json:"code"`
	Status          string    `json:"status"`
	Runtime         int       `json:"runtime"`
	MemoryUsed      int       `json:"memory_used"`
	TotalTestcases  int       `json:"total_testcases"`
	PassedTestcases int       `json:"passed_testcases"`
	FailedTestcases int       `json:"failed_testcases"`
	CreatedAt       time.Time `json:"created_at"`
}

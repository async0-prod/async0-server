package models

import (
	"time"

	"github.com/google/uuid"
)

type Solution struct {
	ID              uuid.UUID `json:"id"`
	ProblemID       uuid.UUID `json:"problem_id"`
	Title           string    `json:"title"`
	Hint            string    `json:"hint"`
	Description     string    `json:"description"`
	Code            string    `json:"code"`
	CodeExplanation string    `json:"code_explanation"`
	Notes           string    `json:"notes"`
	TimeComplexity  string    `json:"time_complexity"`
	SpaceComplexity string    `json:"space_complexity"`
	DifficultyLevel string    `json:"difficulty_level"`
	DisplayOrder    int       `json:"display_order"`
	Author          string    `json:"author"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SolutionBasic struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Hint            string    `json:"hint"`
	Description     string    `json:"description"`
	Code            string    `json:"code"`
	CodeExplanation string    `json:"code_explanation"`
	Notes           string    `json:"notes"`
	TimeComplexity  string    `json:"time_complexity"`
	SpaceComplexity string    `json:"space_complexity"`
	DifficultyLevel string    `json:"difficulty_level"`
	DisplayOrder    int       `json:"display_order"`
	Author          string    `json:"author"`
	IsActive        bool      `json:"is_active"`
}

package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type PostgresSolutionStore struct {
	DB *sql.DB
}

func NewPostgresSolutionStore(db *sql.DB) *PostgresSolutionStore {
	return &PostgresSolutionStore{
		DB: db,
	}
}

type SolutionStore interface {
	GetSolutionsByProblemID(problemID uuid.UUID) ([]models.SolutionBasic, error)
}

func (ps *PostgresSolutionStore) GetSolutionsByProblemID(problemID uuid.UUID) ([]models.SolutionBasic, error) {

	query := `
		SELECT
			id,
			title,
			hint,
			description,
			code,
			code_explanation,
			notes,
			time_complexity,
			space_complexity,
			difficulty_level,
			display_order,
			author,
			is_active
		FROM solutions
		WHERE problem_id = $1 AND is_active = TRUE
	`

	rows, err := ps.DB.Query(query, problemID)
	if err != nil {
		return nil, fmt.Errorf("error querying solutions: %w", err)
	}

	defer rows.Close()

	solutions := []models.SolutionBasic{}
	for rows.Next() {
		var sol models.SolutionBasic
		err := rows.Scan(
			&sol.ID,
			&sol.Title,
			&sol.Hint,
			&sol.Description,
			&sol.Code,
			&sol.CodeExplanation,
			&sol.Notes,
			&sol.TimeComplexity,
			&sol.SpaceComplexity,
			&sol.DifficultyLevel,
			&sol.DisplayOrder,
			&sol.Author,
			&sol.IsActive,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning solution: %w", err)
		}

		solutions = append(solutions, sol)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating solution: %w", err)
	}

	return solutions, nil
}

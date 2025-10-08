package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type AdminPostgresSolutionStore struct {
	DB *sql.DB
}

func NewPostgresAdminSolutionStore(db *sql.DB) *AdminPostgresSolutionStore {
	return &AdminPostgresSolutionStore{
		DB: db,
	}
}

type AdminSolutionStore interface {
	GetSolutionsByProblemID(problemID uuid.UUID) ([]models.SolutionBasic, error)
}

func (ap *AdminPostgresSolutionStore) GetSolutionsByProblemID(problemID uuid.UUID) ([]models.SolutionBasic, error) {
	query := `
		SELECT
			s.id,
			s.title,
			s.hint,
			s.description,
			s.code,
			s.code_explanation,
			s.notes,
			s.time_complexity,
			s.space_complexity,
			s.difficulty_level,
			s.display_order,
			s.author
		FROM solutions s
		WHERE s.problem_id = $1 AND s.is_active = TRUE
	`

	rows, err := ap.DB.Query(query, problemID)
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
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning solution: %w", err)
		}

		solutions = append(solutions, sol)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating solutions: %w", err)
	}

	return solutions, nil
}

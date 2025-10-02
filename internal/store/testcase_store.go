package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type PostgresTestcaseStore struct {
	DB *sql.DB
}

func NewPostgresTestcaseStore(db *sql.DB) *PostgresTestcaseStore {
	return &PostgresTestcaseStore{
		DB: db,
	}
}

type TestcaseStore interface {
	GetTestcasesByProblemID(problemID uuid.UUID) ([]models.Testcase, error)
}

func (ps *PostgresTestcaseStore) GetTestcasesByProblemID(problemID uuid.UUID) ([]models.Testcase, error) {
	query := `
		SELECT
			id,
			problem_id,
			ui,
			input,
			output,
			position,
			is_active,
			created_at
		FROM testcases
		WHERE problem_id = $1
	`

	rows, err := ps.DB.Query(query, problemID)
	if err != nil {
		return nil, fmt.Errorf("error querying testcases: %w", err)
	}

	defer rows.Close()

	testcases := []models.Testcase{}
	for rows.Next() {
		var tc models.Testcase
		err := rows.Scan(
			&tc.ID,
			&tc.ProblemID,
			&tc.UI,
			&tc.Input,
			&tc.Output,
			&tc.Position,
			&tc.IsActive,
			&tc.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning testcase: %w", err)
		}

		testcases = append(testcases, tc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating testcases: %w", err)
	}

	return testcases, nil
}

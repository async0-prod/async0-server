package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type AdminPostgresTestcaseStore struct {
	DB *sql.DB
}

func NewPostgresAdminTestcaseStore(db *sql.DB) *AdminPostgresTestcaseStore {
	return &AdminPostgresTestcaseStore{
		DB: db,
	}
}

type AdminTestcaseStore interface {
	GetTestcasesByProblemID(problemID uuid.UUID) ([]models.TestcaseBasic, error)
}

func (ap *AdminPostgresTestcaseStore) GetTestcasesByProblemID(problemID uuid.UUID) ([]models.TestcaseBasic, error) {

	query := `
		SELECT
			id,
			ui,
			input,
			output
		FROM testcases
		WHERE problem_id = $1
	`
	rows, err := ap.DB.Query(query, problemID)
	if err != nil {
		return nil, fmt.Errorf("error querying testcases: %w", err)
	}

	defer rows.Close()

	testcases := []models.TestcaseBasic{}
	for rows.Next() {
		var tc models.TestcaseBasic
		err := rows.Scan(
			&tc.ID,
			&tc.UI,
			&tc.Input,
			&tc.Output,
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

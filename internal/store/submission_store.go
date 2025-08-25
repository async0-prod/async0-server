package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type PostgresSubmissionStore struct {
	DB *sql.DB
}

func NewPostgresSubmissionStore(db *sql.DB) *PostgresSubmissionStore {
	return &PostgresSubmissionStore{
		DB: db,
	}
}

type SubmissionStore interface {
	CreateSubmission(userID uuid.UUID, problemID uuid.UUID, code string, result models.SubmitSubmissionResponse) error
}

func (ps *PostgresSubmissionStore) CreateSubmission(userID uuid.UUID, problemID uuid.UUID, code string, result models.SubmitSubmissionResponse) error {
	query := `
		INSERT INTO submissions (user_id, problem_id, code, status, total_testcases, passed_testcases, failed_testcases)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// var total_memory int
	// var total_time float64

	// for _, testcaseResult := range result.TestcasesResults {
	// 	total_memory += testcaseResult.TCMemory
	// 	total_time += testcaseResult.TCTime
	// }

	_, err := ps.DB.Exec(query, userID, problemID, code, result.OverallStatus, result.TotalTestcases, result.PassedTestcases, result.TotalTestcases-result.PassedTestcases)
	if err != nil {
		return fmt.Errorf("error running create submission query: %w", err)
	}
	return nil

}

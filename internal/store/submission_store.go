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
	GetSubmissionsByProblemID(userID uuid.UUID, problemID uuid.UUID) ([]models.Submission, error)
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

func (ps *PostgresSubmissionStore) GetSubmissionsByProblemID(userID uuid.UUID, problemID uuid.UUID) ([]models.Submission, error) {
	var submissions []models.Submission

	query := `
		SELECT * FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY created_at DESC
	`

	rows, err := ps.DB.Query(query, userID, problemID)
	if err != nil {
		return nil, fmt.Errorf("error running get submissions by problem id query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		submission := models.Submission{}
		err = rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ProblemID,
			&submission.Code,
			&submission.Status,
			&submission.Runtime,
			&submission.MemoryUsed,
			&submission.TotalTestcases,
			&submission.PassedTestcases,
			&submission.FailedTestcases,
			&submission.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

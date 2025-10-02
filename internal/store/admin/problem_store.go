package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type AdminPostgresProblemStore struct {
	DB *sql.DB
}

func NewPostgresAdminProblemStore(db *sql.DB) *AdminPostgresProblemStore {
	return &AdminPostgresProblemStore{
		DB: db,
	}
}

type AdminProblemStore interface {
	GetAllProblems() ([]models.Problem, error)
	CreateProblem(models.Problem, []uuid.UUID, []uuid.UUID, []models.Testcase, []models.Solution) error
}

func (ap *AdminPostgresProblemStore) GetAllProblems() ([]models.Problem, error) {
	problems := []models.Problem{}

	query := `
		SELECT id, name, slug, link, problem_number, difficulty, starter_code, time_limit, memory_limit, acceptance_rate, total_submissions, successful_submissions, is_active
		FROM problems
	`

	rows, err := ap.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		problem := models.Problem{}
		err := rows.Scan(&problem.ID, &problem.Name, &problem.Slug, &problem.Link, &problem.ProblemNumber, &problem.Difficulty, &problem.StarterCode, &problem.TimeLimit, &problem.MemoryLimit, &problem.AcceptanceRate, &problem.TotalSubmissions, &problem.SuccessfulSubmissions, &problem.IsActive)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

func (ap *AdminPostgresProblemStore) CreateProblem(problem models.Problem, listIDs []uuid.UUID, topicIDs []uuid.UUID, testcases []models.Testcase, solutions []models.Solution) error {

	tx, err := ap.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
			fmt.Printf("rollback error: %v", rErr)
		}
	}()

	// insert problem
	var problemID uuid.UUID
	query := `
		INSERT INTO problems (name, slug, description, link, difficulty, starter_code, time_limit, memory_limit, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
		`
	err = tx.QueryRow(query, problem.Name, problem.Slug, problem.Description, problem.Link, problem.Difficulty, problem.StarterCode, problem.TimeLimit, problem.MemoryLimit, problem.IsActive).Scan(&problemID)
	if err != nil {
		return fmt.Errorf("failed to insert problem: %w", err)
	}

	// insert into problem_topics
	for _, topicID := range topicIDs {
		query := `
			INSERT INTO problem_topics (problem_id, topic_id)
			VALUES ($1, $2)
			`
		_, err := tx.Exec(query, problemID, topicID)
		if err != nil {
			return fmt.Errorf("failed to insert problem_topics: %w", err)
		}
	}

	// insert into problem_lists
	for _, listID := range listIDs {
		query := `
			INSERT INTO list_problems (problem_id, list_id, position)
			VALUES ($1, $2, $3)
			`
		_, err := tx.Exec(query, problemID, listID, problem.ProblemNumber)
		if err != nil {
			return fmt.Errorf("failed to insert problem_lists: %w", err)
		}
	}

	// insert into testcases
	for _, testcase := range testcases {
		query := `
			INSERT INTO testcases (problem_id, ui, input, output, position)
			VALUES ($1, $2, $3, $4, $5)
			`
		_, err := tx.Exec(query, problemID, testcase.UI, testcase.Input, testcase.Output, testcase.Position)
		if err != nil {
			return fmt.Errorf("failed to insert testcases: %w", err)
		}

	}

	// insert solutions only if there is at least one

	if len(solutions) > 0 {
		for _, solution := range solutions {
			query := `
				INSERT INTO solutions (problem_id, title, hint, description, code, code_explanation, notes, time_complexity, space_complexity, difficulty_level, display_order, author, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
				`
			_, err := tx.Exec(query, problemID, solution.Title, solution.Hint, solution.Description, solution.Code, solution.CodeExplanation, solution.Notes, solution.TimeComplexity, solution.SpaceComplexity, solution.DifficultyLevel, solution.DisplayOrder, solution.Author, solution.IsActive)
			if err != nil {
				return fmt.Errorf("failed to insert solutions: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

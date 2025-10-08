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
	GetProblemByID(uuid.UUID) (models.Problem, error)
	UpdateProblem(uuid.UUID, models.Problem, []uuid.UUID, []uuid.UUID, []models.Testcase, []models.Solution) error
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

func (ap *AdminPostgresProblemStore) UpdateProblem(problemID uuid.UUID, problem models.Problem, listIDs []uuid.UUID, topicIDs []uuid.UUID, testcases []models.Testcase, solutions []models.Solution) error {
	tx, err := ap.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if rErr := tx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
			fmt.Printf("rollback error: %v", rErr)
		}
	}()

	// 1️⃣ Update main problem
	query := `
		UPDATE problems
		SET name = $1,
			slug = $2,
			description = $3,
			link = $4,
			difficulty = $5,
			starter_code = $6,
			time_limit = $7,
			memory_limit = $8,
			is_active = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
	`
	_, err = tx.Exec(query,
		problem.Name, problem.Slug, problem.Description, problem.Link,
		problem.Difficulty, problem.StarterCode, problem.TimeLimit, problem.MemoryLimit,
		problem.IsActive, problemID)
	if err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}

	// 2️⃣ Replace topics
	_, err = tx.Exec(`DELETE FROM problem_topics WHERE problem_id = $1`, problemID)
	if err != nil {
		return fmt.Errorf("failed to clear problem_topics: %w", err)
	}
	for _, tid := range topicIDs {
		_, err = tx.Exec(`INSERT INTO problem_topics (problem_id, topic_id) VALUES ($1, $2)`, problemID, tid)
		if err != nil {
			return fmt.Errorf("failed to insert problem_topics: %w", err)
		}
	}

	// 3️⃣ Replace lists
	_, err = tx.Exec(`DELETE FROM list_problems WHERE problem_id = $1`, problemID)
	if err != nil {
		return fmt.Errorf("failed to clear list_problems: %w", err)
	}
	for _, listID := range listIDs {
		_, err = tx.Exec(`INSERT INTO list_problems (problem_id, list_id, position) VALUES ($1, $2, $3)`,
			problemID, listID, problem.ProblemNumber)
		if err != nil {
			return fmt.Errorf("failed to insert list_problems: %w", err)
		}
	}

	// 4️⃣ Replace testcases
	_, err = tx.Exec(`DELETE FROM testcases WHERE problem_id = $1`, problemID)
	if err != nil {
		return fmt.Errorf("failed to clear testcases: %w", err)
	}
	for _, tc := range testcases {
		_, err = tx.Exec(`INSERT INTO testcases (problem_id, ui, input, output, position) VALUES ($1, $2, $3, $4, $5)`,
			problemID, tc.UI, tc.Input, tc.Output, tc.Position)
		if err != nil {
			return fmt.Errorf("failed to insert testcases: %w", err)
		}
	}

	// 5️⃣ Replace solutions
	_, err = tx.Exec(`DELETE FROM solutions WHERE problem_id = $1`, problemID)
	if err != nil {
		return fmt.Errorf("failed to clear solutions: %w", err)
	}
	for _, s := range solutions {
		_, err = tx.Exec(`
			INSERT INTO solutions (problem_id, title, hint, description, code, code_explanation, notes, time_complexity, space_complexity, difficulty_level, display_order, author, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`,
			problemID, s.Title, s.Hint, s.Description, s.Code, s.CodeExplanation,
			s.Notes, s.TimeComplexity, s.SpaceComplexity, s.DifficultyLevel,
			s.DisplayOrder, s.Author, s.IsActive)
		if err != nil {
			return fmt.Errorf("failed to insert solutions: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit update: %w", err)
	}
	return nil
}

func (ap *AdminPostgresProblemStore) GetProblemByID(problemID uuid.UUID) (models.Problem, error) {

	query := `
		SELECT id, name, slug, description, link, problem_number, difficulty, starter_code, time_limit, memory_limit, acceptance_rate, total_submissions, successful_submissions, is_active
		FROM problems
		WHERE id = $1
	`

	row := ap.DB.QueryRow(query, problemID)

	problem := models.Problem{}
	err := row.Scan(&problem.ID, &problem.Name, &problem.Slug, &problem.Description, &problem.Link, &problem.ProblemNumber, &problem.Difficulty, &problem.StarterCode, &problem.TimeLimit, &problem.MemoryLimit, &problem.AcceptanceRate, &problem.TotalSubmissions, &problem.SuccessfulSubmissions, &problem.IsActive)
	if err != nil {
		return models.Problem{}, fmt.Errorf("error running get problem by id query: %w", err)
	}

	if err == sql.ErrNoRows {
		return models.Problem{}, fmt.Errorf("problem not found")
	}

	return problem, nil

}

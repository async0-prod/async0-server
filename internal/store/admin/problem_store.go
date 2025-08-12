package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

// {
//     "name": "Two Sum",
//     "slug": "two-sum",
//     "link": "asaas",
//     "difficulty": "Easy",
//     "starter_code": {},
//     "solution_code": {},
//     "time_limit": 2000,
//     "memory_limit": 256,
//     "is_active": true,
//     "topics": [
//         "2ce5b4f9-7ebb-4ce8-9787-558c11ca86ad",
//         "660b2f88-7b8f-4009-8cb5-184c45780cdf"
//     ],
//     "lists": [
//         "daa215cc-c573-466e-94ec-902ec072c9f7",
//         "7ab83b2c-66e8-4b00-9858-f340212e481a"
//     ],
//     "testCases": [
//         {
//             "name": "asas",
//             "input": "asas",
//             "output": "asas"
//         }
//     ]
// }

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
	CreateProblem(models.Problem, []uuid.UUID, []uuid.UUID, []models.TestCase) error
}

func (ap *AdminPostgresProblemStore) GetAllProblems() ([]models.Problem, error) {
	problems := []models.Problem{}

	query := `
		SELECT id, name, slug, link, problem_number, difficulty, starter_code, solution_code, time_limit, memory_limit, acceptance_rate, total_submissions, successful_submissions, is_active
		FROM problems
	`

	rows, err := ap.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		problem := models.Problem{}
		err := rows.Scan(&problem.ID, &problem.Name, &problem.Slug, &problem.Link, &problem.ProblemNumber, &problem.Difficulty, &problem.StarterCode, &problem.SolutionCode, &problem.TimeLimit, &problem.MemoryLimit, &problem.AcceptanceRate, &problem.TotalSubmissions, &problem.SuccessfulSubmissions, &problem.IsActive)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

func (ap *AdminPostgresProblemStore) CreateProblem(problem models.Problem, listIDs []uuid.UUID, topicIDs []uuid.UUID, testcases []models.TestCase) error {

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
	query := `
		INSERT INTO problems (name, slug, link, difficulty, starter_code, solution_code, time_limit, memory_limit, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
		`
	problemID, err := tx.Exec(query, problem.Name, problem.Slug, problem.Link, problem.Difficulty, problem.StarterCode, problem.SolutionCode, problem.TimeLimit, problem.MemoryLimit, problem.IsActive)
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
			INSERT INTO problem_lists (problem_id, list_id)
			VALUES ($1, $2)
			`
		_, err := tx.Exec(query, problemID, listID)
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

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

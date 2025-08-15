package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

var ErrProblemNotFound = errors.New("problem not found")

type PostgresProblemStore struct {
	DB *sql.DB
}

func NewPostgresProblemStore(db *sql.DB) *PostgresProblemStore {
	return &PostgresProblemStore{DB: db}
}

type ProblemStore interface {
	GetProblemBySlug(slug string) (*models.Problem, error)
	GetListProblems(listID uuid.UUID) ([]models.Problem, error)
}

func (p *PostgresProblemStore) GetProblemBySlug(slug string) (*models.Problem, error) {
	query := `
		SELECT id, name, slug, link, problem_number, difficulty, starter_code, time_limit, memory_limit, acceptance_rate, total_submissions, successful_submissions, is_active
		FROM problems
		WHERE slug = $1
	`

	var problem models.Problem
	err := p.DB.QueryRow(query, slug).Scan(
		&problem.ID,
		&problem.Name,
		&problem.Slug,
		&problem.Link,
		&problem.ProblemNumber,
		&problem.Difficulty,
		&problem.StarterCode,
		&problem.TimeLimit,
		&problem.MemoryLimit,
		&problem.AcceptanceRate,
		&problem.TotalSubmissions,
		&problem.SuccessfulSubmissions,
		&problem.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, ErrProblemNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error running get problem by slug query: %w", err)
	}

	return &problem, nil

}

func (pg *PostgresProblemStore) GetListProblems(listID uuid.UUID) ([]models.Problem, error) {
	query := `
		SELECT id, name, slug, link, problem_number, difficulty, starter_code, time_limit, memory_limit, acceptance_rate, total_submissions, successful_submissions, is_active
		FROM problems
		WHERE id IN (
			SELECT problem_id
			FROM list_problems
			WHERE list_id = $1
		)
		`

	result, err := pg.DB.Query(query, listID)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	problems := []models.Problem{}

	for result.Next() {
		problem := models.Problem{}
		err := result.Scan(&problem.ID, &problem.Name, &problem.Slug, &problem.Link, &problem.ProblemNumber, &problem.Difficulty, &problem.StarterCode, &problem.TimeLimit, &problem.MemoryLimit, &problem.AcceptanceRate, &problem.TotalSubmissions, &problem.SuccessfulSubmissions, &problem.IsActive)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

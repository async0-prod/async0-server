package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

type TanstackTableProblem struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Difficulty string    `json:"difficulty"`
	ListNames  []string  `json:"list_names"`
	TopicNames []string  `json:"topic_names"`
	HasSolved  *bool     `json:"has_solved"`
}

type ProblemStore interface {
	GetProblemBySlug(slug string) (*models.Problem, error)
	GetTanstackTableProblemsByListID(userID *uuid.UUID, listID uuid.UUID) ([]TanstackTableProblem, error)
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

func (pg *PostgresProblemStore) GetTanstackTableProblemsByListID(userID *uuid.UUID, listID uuid.UUID) ([]TanstackTableProblem, error) {

	query := `
	SELECT
		p.id,
		p.name,
		p.slug,
		p.difficulty,
		STRING_AGG(DISTINCT l.name, ', ') as list_names,
		STRING_AGG(DISTINCT t.name, ', ') as topic_names,
		CASE
			WHEN $1::UUID IS NULL THEN NULL
			WHEN s.user_id IS NOT NULL THEN true
			ELSE false
		END as has_solved
	FROM problems p
	INNER JOIN list_problems lp_filter ON p.id = lp_filter.problem_id AND lp_filter.list_id = $2
	LEFT JOIN list_problems lp ON p.id = lp.problem_id
	LEFT JOIN lists l ON lp.list_id = l.id AND l.is_active = true
	LEFT JOIN problem_topics pt ON p.id = pt.problem_id
	LEFT JOIN topics t ON pt.topic_id = t.id AND t.is_active = true
	LEFT JOIN (
		SELECT DISTINCT user_id, problem_id
		FROM submissions
		WHERE ($1::UUID IS NULL OR user_id = $1) AND status = 'AC'
	) s ON p.id = s.problem_id AND $1::UUID IS NOT NULL
	WHERE p.is_active = true
	GROUP BY p.id, p.name, p.slug, p.difficulty, s.user_id, lp_filter.position
	ORDER BY lp_filter.position, p.problem_number;
	`

	rows, err := pg.DB.Query(query, userID, listID)
	if err != nil {
		return nil, fmt.Errorf("error running get tanstack table query: %w", err)
	}
	defer rows.Close()

	var tableProblems []TanstackTableProblem

	for rows.Next() {
		var tableProblem TanstackTableProblem
		var listNamesStr, topicNamesStr sql.NullString

		err := rows.Scan(
			&tableProblem.ID,
			&tableProblem.Name,
			&tableProblem.Slug,
			&tableProblem.Difficulty,
			&listNamesStr,
			&topicNamesStr,
			&tableProblem.HasSolved,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if listNamesStr.Valid && listNamesStr.String != "" {
			tableProblem.ListNames = strings.Split(listNamesStr.String, ", ")
		} else {
			tableProblem.ListNames = []string{}
		}

		if topicNamesStr.Valid && topicNamesStr.String != "" {
			tableProblem.TopicNames = strings.Split(topicNamesStr.String, ", ")
		} else {
			tableProblem.TopicNames = []string{}
		}

		tableProblems = append(tableProblems, tableProblem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tableProblems, nil
}

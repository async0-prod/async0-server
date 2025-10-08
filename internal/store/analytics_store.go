package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type PostgresAnalyticsStore struct {
	db *sql.DB
}

func NewPostgresAnalyticsStore(db *sql.DB) *PostgresAnalyticsStore {
	return &PostgresAnalyticsStore{
		db: db,
	}
}

type AnalyticsStore interface {
	GetCardAnalyticsByListID(userID uuid.UUID, listID uuid.UUID) (models.CardAnalytics, error)
	GetCardAnalyticsByListIDNoUser(listID uuid.UUID) (models.CardAnalytics, error)
}

func (pg *PostgresAnalyticsStore) GetCardAnalyticsByListID(userID uuid.UUID, listID uuid.UUID) (models.CardAnalytics, error) {

	query := `
	SELECT
		l.name,
		COUNT(DISTINCT lp.problem_id) as total_questions,
		COUNT(DISTINCT CASE WHEN s_accepted.problem_id IS NOT NULL THEN lp.problem_id END) as total_solved,
		COUNT(DISTINCT sol.id) as total_solutions,
		COUNT(s_all.id) as total_user_attempts,

		COUNT(DISTINCT CASE WHEN p.difficulty = 'EASY' THEN lp.problem_id END) as total_easy_q,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'MEDIUM' THEN lp.problem_id END) as total_medium_q,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'HARD' THEN lp.problem_id END) as total_hard_q,

		COUNT(DISTINCT CASE WHEN p.difficulty = 'EASY' AND s_accepted.problem_id IS NOT NULL THEN lp.problem_id END) as total_easy_solved,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'MEDIUM' AND s_accepted.problem_id IS NOT NULL THEN lp.problem_id END) as total_medium_solved,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'HARD' AND s_accepted.problem_id IS NOT NULL THEN lp.problem_id END) as total_hard_solved

	FROM lists l
	LEFT JOIN list_problems lp ON l.id = lp.list_id
	LEFT JOIN problems p ON lp.problem_id = p.id AND p.is_active = true
	LEFT JOIN solutions sol ON p.id = sol.problem_id AND sol.is_active = true
	LEFT JOIN submissions s_all ON p.id = s_all.problem_id AND s_all.user_id = $1
	LEFT JOIN submissions s_accepted ON p.id = s_accepted.problem_id AND s_accepted.user_id = $1 AND s_accepted.status = 'AC'
	WHERE l.id = $2 AND l.is_active = true
	GROUP BY l.id, l.name;
	`

	var analytics models.CardAnalytics
	err := pg.db.QueryRow(query, userID, listID).Scan(
		&analytics.Name,
		&analytics.TotalQuestions,
		&analytics.TotalSolved,
		&analytics.TotalSolutions,
		&analytics.TotalUserAttempts,
		&analytics.TotalEasyQ,
		&analytics.TotalMediumQ,
		&analytics.TotalHardQ,
		&analytics.TotalEasySolved,
		&analytics.TotalMediumSolved,
		&analytics.TotalHardSolved,
	)

	if err == sql.ErrNoRows {
		return models.CardAnalytics{}, fmt.Errorf("card analytics not found")
	}

	if err != nil {
		return models.CardAnalytics{}, fmt.Errorf("error running get card analytics query: %w", err)
	}

	return analytics, nil
}

func (pg *PostgresAnalyticsStore) GetCardAnalyticsByListIDNoUser(listID uuid.UUID) (models.CardAnalytics, error) {
	query := `
	SELECT
		l.name,
		COUNT(DISTINCT lp.problem_id) as total_questions,
		COUNT(DISTINCT sol.id) as total_solutions,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'EASY' THEN lp.problem_id END) as total_easy_q,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'MEDIUM' THEN lp.problem_id END) as total_medium_q,
		COUNT(DISTINCT CASE WHEN p.difficulty = 'HARD' THEN lp.problem_id END) as total_hard_q
	FROM lists l
	LEFT JOIN list_problems lp ON l.id = lp.list_id
	LEFT JOIN problems p ON lp.problem_id = p.id AND p.is_active = true
	LEFT JOIN solutions sol ON p.id = sol.problem_id AND sol.is_active = true
	WHERE l.id = $1 AND l.is_active = true
	GROUP BY l.id, l.name;
	`

	var analytics models.CardAnalytics
	err := pg.db.QueryRow(query, listID).Scan(
		&analytics.Name,
		&analytics.TotalQuestions,
		&analytics.TotalSolutions,
		&analytics.TotalEasyQ,
		&analytics.TotalMediumQ,
		&analytics.TotalHardQ,
	)

	if err == sql.ErrNoRows {
		return models.CardAnalytics{}, fmt.Errorf("card analytics not found")
	}

	if err != nil {
		return models.CardAnalytics{}, fmt.Errorf("error running get card analytics query: %w", err)
	}

	analytics.TotalSolved = 0
	analytics.TotalUserAttempts = 0
	analytics.TotalEasySolved = 0
	analytics.TotalMediumSolved = 0
	analytics.TotalHardSolved = 0

	return analytics, nil
}

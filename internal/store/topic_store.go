package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type PostgresTopicStore struct {
	DB *sql.DB
}

func NewPostgresTopicStore(db *sql.DB) *PostgresTopicStore {
	return &PostgresTopicStore{
		DB: db,
	}
}

type StubbedProblem struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Difficulty string    `json:"difficulty"`
	Position   int       `json:"position"`
}

type TopicsWithProblems struct {
	Name     string           `json:"name"`
	Problems []StubbedProblem `json:"problems"`
}

type TopicStore interface {
	GetAllTopicsByListID(listID uuid.UUID) ([]models.Topic, error)
	GetAllTopicsAndProblemsByListID(listID uuid.UUID) ([]TopicsWithProblems, error)
	GetAllTopics() ([]models.Topic, error)
}

func (ps *PostgresTopicStore) GetAllTopicsByListID(listID uuid.UUID) ([]models.Topic, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			is_active,
			display_order,
			created_at,
			updated_at
		FROM topics
		WHERE list_id = $1
	`

	rows, err := ps.DB.Query(query, listID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	topics := []models.Topic{}
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Slug,
			&topic.IsActive,
			&topic.DisplayOrder,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}

func (ps *PostgresTopicStore) GetAllTopicsAndProblemsByListID(listID uuid.UUID) ([]TopicsWithProblems, error) {
	query := `
		WITH ordered_topics AS (
		    SELECT DISTINCT
		        t.id,
		        t.name,
		        t.display_order
		    FROM topics t
		    JOIN problem_topics pt ON t.id = pt.topic_id
		    JOIN problems p ON pt.problem_id = p.id
		    JOIN list_problems lp ON p.id = lp.problem_id
		    WHERE lp.list_id = $1
		      AND t.is_active = true
		      AND p.is_active = true
		    ORDER BY t.display_order
		)
		SELECT
		    ot.name AS topic_name,
		    COALESCE(
		        json_agg(
		            json_build_object(
		                'id', p.id,
		                'name', p.name,
		                'slug', p.slug,
		                'difficulty', p.difficulty,
		                'position', lp.position
		            ) ORDER BY lp.position
		        ) FILTER (WHERE p.id IS NOT NULL),
		        '[]'::json
		    ) AS problems
		FROM ordered_topics ot
		LEFT JOIN problem_topics pt ON ot.id = pt.topic_id
		LEFT JOIN problems p ON pt.problem_id = p.id AND p.is_active = true
		LEFT JOIN list_problems lp ON p.id = lp.problem_id AND lp.list_id = $1
		GROUP BY ot.name, ot.display_order
		ORDER BY ot.display_order;
	`

	rows, err := ps.DB.Query(query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var topics []TopicsWithProblems
	for rows.Next() {
		var (
			topicName    string
			problemsJSON string
		)

		err := rows.Scan(&topicName, &problemsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var problems []StubbedProblem
		err = json.Unmarshal([]byte(problemsJSON), &problems)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal problems JSON: %w", err)
		}

		topics = append(topics, TopicsWithProblems{
			Name:     topicName,
			Problems: problems,
		})
	}

	return topics, nil
}

func (ps *PostgresTopicStore) GetAllTopics() ([]models.Topic, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			is_active,
			display_order,
			created_at,
			updated_at
		FROM topics
	`

	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get all topics: %w", err)
	}

	defer rows.Close()

	topics := []models.Topic{}
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Slug,
			&topic.IsActive,
			&topic.DisplayOrder,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

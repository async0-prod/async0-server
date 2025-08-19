package store

import (
	"database/sql"
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
	Topic    string           `json:"topic"`
	Problems []StubbedProblem `json:"problems"`
}

type TopicStore interface {
	GetAllTopicsByListID(listID uuid.UUID) ([]models.Topic, error)
	GetAllTopicsAndProblemsByListID(listID uuid.UUID) ([]TopicsWithProblems, error)
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
		SELECT
		    t.name AS topic,
		    p.id,
		    p.name,
		    p.slug,
		    p.difficulty,
		    lp.position
		FROM list_problems lp
		JOIN problems p ON lp.problem_id = p.id
		JOIN problem_topics pt ON p.id = pt.problem_id
		JOIN topics t ON pt.topic_id = t.id
		WHERE lp.list_id = $1
		ORDER BY t.display_order, lp.position;
	`

	rows, err := ps.DB.Query(query, listID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	grouped := make(map[string][]StubbedProblem)

	for rows.Next() {
		var (
			topicName string
			problem   StubbedProblem
		)

		err := rows.Scan(
			&topicName,
			&problem.ID,
			&problem.Name,
			&problem.Slug,
			&problem.Difficulty,
			&problem.Position,
		)

		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		grouped[topicName] = append(grouped[topicName], problem)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	var topics []TopicsWithProblems
	for topic, problems := range grouped {
		topics = append(topics, TopicsWithProblems{
			Topic:    topic,
			Problems: problems,
		})
	}

	return topics, nil
}

package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type AdminPostgresTopicStore struct {
	DB *sql.DB
}

func NewPostgresAdminTopicStore(db *sql.DB) *AdminPostgresTopicStore {
	return &AdminPostgresTopicStore{
		DB: db,
	}
}

type AdminTopicStore interface {
	GetAllTopics() ([]models.Topic, error)
	GetTopicsByProblemID(problemID uuid.UUID) ([]models.TopicBasic, error)
}

func (a *AdminPostgresTopicStore) GetAllTopics() ([]models.Topic, error) {
	topics := []models.Topic{}

	query := `
		SELECT id, name, slug, is_active, display_order
		FROM topics
	`

	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		topic := models.Topic{}
		err := rows.Scan(&topic.ID, &topic.Name, &topic.Slug, &topic.IsActive, &topic.DisplayOrder)
		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (a *AdminPostgresTopicStore) GetTopicsByProblemID(problemID uuid.UUID) ([]models.TopicBasic, error) {

	query := `
		SELECT
			t.id,
			t.name
		FROM problem_topics pt
		JOIN topics t ON pt.topic_id = t.id
		WHERE pt.problem_id = $1;
	`

	rows, err := a.DB.Query(query, problemID)
	if err != nil {
		return nil, fmt.Errorf("error running get topics by problem id query: %w", err)
	}

	defer rows.Close()

	topics := []models.TopicBasic{}
	for rows.Next() {
		topic := models.TopicBasic{}
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning topic: %w", err)
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

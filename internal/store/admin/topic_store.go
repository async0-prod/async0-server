package admin

import (
	"database/sql"

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

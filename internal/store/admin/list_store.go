package admin

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
)

type AdminPostgresListStore struct {
	DB *sql.DB
}

func NewPostgresAdminListStore(db *sql.DB) *AdminPostgresListStore {
	return &AdminPostgresListStore{
		DB: db,
	}
}

type AdminListStore interface {
	GetAllLists() ([]models.List, error)
	GetListsByProblemID(problemID uuid.UUID) ([]models.ListBasic, error)
}

func (ap *AdminPostgresListStore) GetAllLists() ([]models.List, error) {
	query := `
		SELECT id, name, slug, total_problems, is_active, display_order, created_at, updated_at
		FROM lists
	`

	result, err := ap.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error running get all lists query: %w", err)
	}

	defer result.Close()

	lists := []models.List{}

	for result.Next() {
		list := models.List{}
		err := result.Scan(&list.ID, &list.Name, &list.Slug, &list.TotalProblems, &list.IsActive, &list.DisplayOrder, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning list: %w", err)
		}

		lists = append(lists, list)
	}

	return lists, nil
}

func (ap *AdminPostgresListStore) GetListsByProblemID(problemID uuid.UUID) ([]models.ListBasic, error) {

	query := `
		SELECT
			l.id,
			l.name
		FROM list_problems lp
		JOIN lists l ON lp.list_id = l.id
		WHERE lp.problem_id = $1;
	`

	result, err := ap.DB.Query(query, problemID)
	if err != nil {
		return nil, fmt.Errorf("error running get lists by problem id query: %w", err)
	}

	defer result.Close()

	lists := []models.ListBasic{}
	for result.Next() {
		list := models.ListBasic{}
		err := result.Scan(
			&list.ID,
			&list.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list: %w", err)
		}

		lists = append(lists, list)
	}

	return lists, nil
}

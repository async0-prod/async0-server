package admin

import (
	"database/sql"

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
}

func (ap *AdminPostgresListStore) GetAllLists() ([]models.List, error) {
	query := `
		SELECT id, name, slug, total_problems, is_active, display_order, created_at, updated_at
		FROM lists
	`

	result, err := ap.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	lists := []models.List{}

	for result.Next() {
		list := models.List{}
		err := result.Scan(&list.ID, &list.Name, &list.Slug, &list.TotalProblems, &list.IsActive, &list.DisplayOrder, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}

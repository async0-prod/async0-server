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

func (a *AdminPostgresListStore) GetAllLists() ([]models.List, error) {
	lists := []models.List{}

	query := `
		SELECT id, name, slug, COALESCE(link, '') AS link, COALESCE(author, '') AS author, total_problems, is_active, display_order, created_at, updated_at
		FROM lists
	`

	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		list := models.List{}
		err := rows.Scan(&list.ID, &list.Name, &list.Slug, &list.Link, &list.Author, &list.TotalProblems, &list.IsActive, &list.DisplayOrder, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}

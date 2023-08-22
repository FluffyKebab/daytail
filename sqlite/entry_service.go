package sqlite

import (
	"database/sql"

	"github.com/FluffyKebab/daytail"
)

type EntryService struct {
	DB *sql.DB
}

var _ daytail.EntryService = EntryService{}

func (s EntryService) Entry(id int) error {
	return nil
}

func (s EntryService) UserEntries(userId int) ([]daytail.Entry, error) {
	rows, err := s.DB.Query("SELECT * FROM entries WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []daytail.Entry
	for rows.Next() {
		var entry daytail.Entry
		err := rows.Scan(&entry.Title, &entry.Text)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s EntryService) CreateEntry(entry daytail.Entry) (int, error) {
	var id int
	err := s.DB.QueryRow(
		"INSERT INTO TABLES entries (title, text) VALUES ($1, $2) RETURNING id",
		entry.Title,
		entry.Title,
	).Scan(&id)
	return id, err
}

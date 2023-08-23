package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Migrate() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name string)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userId INTEGER, 
		title string, 
		text string)
	`)

	return db, nil
}

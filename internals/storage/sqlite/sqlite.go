package sqlite

import (
	"database/sql"

	"github.com/3shaan/students-api/internals/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, createTableErr := db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if createTableErr != nil {
		return nil, createTableErr
	}
	return &Sqlite{
		Db: db,
	}, nil

}

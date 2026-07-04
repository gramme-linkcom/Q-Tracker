package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	dbPath := fmt.Sprintf("%s/service.db", "./data")
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	if _, err := database.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

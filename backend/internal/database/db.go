package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct{ *sql.DB }

func NewSQLiteStore(driverName, dataSourceName, schemaFile string) (*SQLiteStore, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Read schema file
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}

	// Execute schema
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{DB: db}, nil
}

package database

import (
	"database/sql"
	"os"
	"testing"
)

func TestNewDBStorage_Success(t *testing.T) {
	schema := "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);"
	f, err := os.CreateTemp("", "schema-*.sql")
	if err != nil {
		t.Fatalf("create temp schema file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(schema); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close schema file: %v", err)
	}

	db, err := NewSQLiteStore("sqlite3", ":memory:", f.Name())
	if err != nil {
		t.Fatalf("NewDBStorage returned error: %v", err)
	}
	if db == nil || db.DB == nil {
		t.Fatalf("expected non-nil DBStorage and DB")
	}
	defer db.DB.Close()

	// Verify we can insert and query the table created by the schema.
	if _, err := db.DB.Exec("INSERT INTO users(name) VALUES(?)", "alice"); err != nil {
		t.Fatalf("insert into users failed: %v", err)
	}
	var count int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		t.Fatalf("query users failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row in users, got %d", count)
	}
}

func TestNewDBStorage_MissingSchemaFile(t *testing.T) {
	_, err := NewSQLiteStore("sqlite3", ":memory:", "does-not-exist.sql")
	if err == nil {
		t.Fatalf("expected error for missing schema file, got nil")
	}
}

func TestNewDBStorage_InvalidSQL(t *testing.T) {
	// Write an invalid SQL schema
	f, err := os.CreateTemp("", "bad-schema-*.sql")
	if err != nil {
		t.Fatalf("create temp schema file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("THIS IS NOT SQL;"); err != nil {
		t.Fatalf("write bad schema: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close schema file: %v", err)
	}

	_, err = NewSQLiteStore("sqlite3", ":memory:", f.Name())
	if err == nil {
		t.Fatalf("expected error for invalid SQL schema, got nil")
	}
}

func TestNewDBStorage_UnknownDriver(t *testing.T) {
	// valid schema file but unknown driver should cause sql.Open to fail
	f, err := os.CreateTemp("", "schema-*.sql")
	if err != nil {
		t.Fatalf("create temp schema file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("CREATE TABLE x(id INTEGER);"); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close schema file: %v", err)
	}

	_, err = NewSQLiteStore("unknown_driver_foobar", "dsn", f.Name())
	if err == nil {
		t.Fatalf("expected error for unknown driver, got nil")
	}
}

// Additional small sanity check to ensure DB returned is usable (optional)
func TestNewDBStorage_DBPing(t *testing.T) {
	schema := "CREATE TABLE pingtest (id INTEGER PRIMARY KEY);"
	f, err := os.CreateTemp("", "schema-*.sql")
	if err != nil {
		t.Fatalf("create temp schema file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(schema); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close schema file: %v", err)
	}

	db, err := NewSQLiteStore("sqlite3", ":memory:", f.Name())
	if err != nil {
		t.Fatalf("NewDBStorage returned error: %v", err)
	}
	defer db.DB.Close()

	if err := db.DB.Ping(); err != nil {
		t.Fatalf("db ping failed: %v", err)
	}

	// ensure DB implements basic query
	if _, err := db.DB.Exec("INSERT INTO pingtest DEFAULT VALUES"); err != nil {
		t.Fatalf("insert into pingtest failed: %v", err)
	}
	var cnt int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM pingtest").Scan(&cnt); err != nil {
		t.Fatalf("query pingtest failed: %v", err)
	}
	if cnt != 1 {
		t.Fatalf("expected 1 row in pingtest, got %d", cnt)
	}
}

// helper to ensure the package builds when sql is unused in tests
var _ = sql.ErrNoRows

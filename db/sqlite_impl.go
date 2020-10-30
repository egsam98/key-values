package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var _ DBTX = (*SQLite)(nil)

const DatabaseName = "db.sqlite"

// Реализация СУБД-интерфейса
type SQLite struct {
	db *sql.DB
}

func NewSQLite() (*SQLite, error) {
	db, err := sql.Open("sqlite3", DatabaseName)
	if err != nil {
		return nil, err
	}
	return &SQLite{db: db}, nil
}

func (s SQLite) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s SQLite) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}

func (s SQLite) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s SQLite) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

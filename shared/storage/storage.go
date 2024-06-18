package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(host string, port int, user string, password string, dbName string) (*PostgresStore, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) BeginTransaction() (*sql.Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *PostgresStore) QueryData(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *PostgresStore) QueryRowData(ctx context.Context, query string, args ...interface{}) *sql.Row {
  return s.db.QueryRowContext(ctx, query, args...)
}

func (s *PostgresStore) ExecData(ctx context.Context, query string, args ...interface{}) error {
  _, err := s.db.ExecContext(ctx, query, args...)
  return err
}

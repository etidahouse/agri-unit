package misc

import (
	"database/sql"
	"fmt"
	"testing"

	go_sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type MockQuerier struct {
	Db   *sql.DB
	Mock go_sqlmock.Sqlmock
}

func NewMockQuerier(t *testing.T) (*MockQuerier, go_sqlmock.Sqlmock, error) {
	db, mock, err := go_sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
		return nil, nil, fmt.Errorf("failed to create sqlmock: %w", err)
	}
	return &MockQuerier{Db: db, Mock: mock}, mock, nil
}

func (m *MockQuerier) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.Db.Query(query, args...)
}

func (m *MockQuerier) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.Db.Exec(query, args...)
}

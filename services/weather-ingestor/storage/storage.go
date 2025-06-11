package storage

import "database/sql"

type DBQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type realDBQuerier struct {
	db *sql.DB
}

func NewRealDBQuerier(db *sql.DB) DBQuerier {
	return &realDBQuerier{db: db}
}

func (r *realDBQuerier) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.Query(query, args...)
}

func (r *realDBQuerier) Exec(query string, args ...interface{}) (sql.Result, error) {
	return r.db.Exec(query, args...)
}

package sql

import "database/sql"

// Sql database abstraction layer
// Source: https://johnjianwang.medium.com/database-abstractions-for-golang-ae252911de6f

type Scannable interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

// Both sql.DB and sql.Tx implement the Database interface. This allows us to use either one in our repository
type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// This scanrows func ensures that the rows are closed after the scanFunc is called
// Must be called every time a query is made
func ScanRows(r Rows, scanFunc func(row Scannable) error) error {
	var closeErr error
	defer func() {
		if err := r.Close(); err != nil {
			closeErr = err
		}
	}()

	var scanErr error
	for r.Next() {
		err := scanFunc(r)
		if err != nil {
			scanErr = err
			break
		}
	}
	if r.Err() != nil {
		return r.Err()
	}
	if scanErr != nil {
		return scanErr
	}

	return closeErr
}

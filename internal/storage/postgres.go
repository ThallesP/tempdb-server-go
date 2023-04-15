package storage

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)


func BootstrapPostgres(uri string, timeout time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ClosePostgres(db *sql.DB) error {
	return db.Close()
}

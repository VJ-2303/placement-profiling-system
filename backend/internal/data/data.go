package data

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

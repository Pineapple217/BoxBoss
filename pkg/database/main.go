package database

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	ddl     string
	queries *Queries
)

func Init(databaseSource string) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", databaseSource)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries = New(db)
}

func GetQueries() *Queries {
	return queries
}

package db

import (
	"context"
	_ "database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresConnection(ctx context.Context, dburl string) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "postgres", dburl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return db
}

func ClosePostgresConnection(db *sqlx.DB) error {
	err := db.Close()
	if err != nil {
		log.Println("error in close connection to db")
	}
	return nil
}

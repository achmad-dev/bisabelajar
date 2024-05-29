package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

func NewPostgresConnection(ctx context.Context, dburl string) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "pgx", dburl)
	if err != nil {
		log.Fatalln("error connect to db")
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

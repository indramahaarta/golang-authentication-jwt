package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v5"
)

var DB *sql.DB

func ConnectDatabase() (*sql.DB, error) {
	d, err := sql.Open("pgx", "host=localhost port=5432 dbname=jwt user=postgres password=user")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = d.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	log.Println("Connected to database")

	DB = d
	return d, nil
}

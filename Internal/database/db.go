package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func connect() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s, port=%s, user=%s, password=%s, name=%s, sslmode=%s",

		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("JWT_SECRET"),
	)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err

	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(25)
	log.Println("db is connected", os.Getenv("DB_NAME"))
	return db, nil

}

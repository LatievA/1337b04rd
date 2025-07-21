package repository

import (
	"1337b04rd/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const POSTGRES = "postgres"

func ConnectToDB(cfg *config.DBConfig) (*sql.DB, error) {
	conn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable")
	db, err := sql.Open(POSTGRES, conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
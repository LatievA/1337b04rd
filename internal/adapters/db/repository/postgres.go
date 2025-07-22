package repository

import (
	"1337b04rd/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const POSTGRES = "postgres"

func ConnectToDB(cfg *config.DBConfig) (*sql.DB, error) {
	conn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", POSTGRES, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open(POSTGRES, conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

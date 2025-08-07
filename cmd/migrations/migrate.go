package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	conn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "hacker", "board", "db", "5432", "1337b04rd")
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Read the init.sql file
	sqlBytes, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalf("Failed to read init.sql: %v", err)
	}

	// Execute SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to execute SQL: %v", err)
	}

	log.Println("Database initialized successfully")
}

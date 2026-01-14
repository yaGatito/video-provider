package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/user-records")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to the database successfully")
	defer db.Close()

	sqlBytes, err := os.ReadFile("/cmd/migrations/001_create_users.sql")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Migration file loaded successfully")

	statements := strings.Split(string(sqlBytes), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err = db.Exec(stmt)
		if err != nil {
			log.Fatalf("Error executing statement: %s\n%v", stmt, err)
		}
	}
	log.Printf("Migration performed successfully")
}

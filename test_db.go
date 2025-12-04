package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// Try different connection strings
	connStrings := []string{
		"postgres://postgres:postgres@localhost:5432/tasks_db?sslmode=disable",
		"postgres://postgres:admin@localhost:5432/tasks_db?sslmode=disable",
		"postgres://postgres:password@localhost:5432/tasks_db?sslmode=disable",
	}

	for _, connStr := range connStrings {
		fmt.Printf("Testing: %s\n", connStr)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			fmt.Printf("  ❌ Open failed: %v\n", err)
			continue
		}

		err = db.Ping()
		if err != nil {
			fmt.Printf("  ❌ Ping failed: %v\n", err)
		} else {
			fmt.Printf("  ✅ SUCCESS! This password works!\n")
			db.Close()
			break
		}
		db.Close()
	}
}

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:vishsql273@localhost:5432/tasks_db?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("database connection failed:%w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Database connected successfully")

	if err = createTables(); err != nil {
		return fmt.Errorf("table creation failed: %w", err)
	}

	return nil

}

func createTables() error {

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY, -- Auto-increment ID
	title VARCHAR(255) NOT NULL, -- Title required hai
	description TEXT, -- Description optional 
	completed BOOLEAN DEFAULT FALSE, -- Default false
	created_at TIMESTAMP DEFAULT NOW(), -- Automatic timestamp
	updated_at TIMESTAMP DEFAULT NOW() 
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println(" Table ready!")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

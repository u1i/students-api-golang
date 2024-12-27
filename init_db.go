package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parse command line flags
	dbPath := flag.String("db", "./students.db", "Path to SQLite database file")
	flag.Parse()

	// Ensure the directory exists
	dbDir := filepath.Dir(*dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create table
	createTableSQL := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		linkedin_profile TEXT,
		phone TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert sample data
	sampleData := []struct {
		name           string
		email          string
		linkedinProfile string
		phone          string
	}{
		{
			name:           "John Doe",
			email:          "john.doe@example.com",
			linkedinProfile: "https://linkedin.com/in/johndoe",
			phone:          "+1-555-123-4567",
		},
		{
			name:           "Jane Smith",
			email:          "jane.smith@example.com",
			linkedinProfile: "https://linkedin.com/in/janesmith",
			phone:          "+1-555-234-5678",
		},
		{
			name:           "Bob Johnson",
			email:          "bob.johnson@example.com",
			linkedinProfile: "https://linkedin.com/in/bobjohnson",
			phone:          "+1-555-345-6789",
		},
	}

	insertSQL := `INSERT OR REPLACE INTO students 
		(name, email, linkedin_profile, phone) 
		VALUES (?, ?, ?, ?)`

	for _, data := range sampleData {
		_, err = db.Exec(insertSQL,
			data.name,
			data.email,
			data.linkedinProfile,
			data.phone,
		)
		if err != nil {
			log.Printf("Error inserting data for %s: %v", data.name, err)
			continue
		}
	}

	log.Printf("Database initialized successfully at: %s", *dbPath)
}

package main

import (
	_ "embed"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed swagger.yaml
var swaggerContent []byte

type Student struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	LinkedinProfile string `json:"linkedin_profile"`
	Phone          string `json:"phone"`
}

type WelcomeResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Swagger     string `json:"swagger_url"`
	Version     string `json:"version"`
}

var db *sql.DB

func initDB(dbPath string) error {
	var err error
	
	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Check if database file exists
	fileExists := false
	if _, err := os.Stat(dbPath); err == nil {
		fileExists = true
		// Check read permission by attempting to open the file
		file, err := os.OpenFile(dbPath, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("database file exists but is not readable: %v", err)
		}
		file.Close()

		// Check write permission
		file, err = os.OpenFile(dbPath, os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("database file exists but is not writable: %v", err)
		}
		file.Close()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check database file: %v", err)
	}

	// Check if directory is writable if file doesn't exist
	if !fileExists {
		// Try to create a temporary file in the directory
		tempFile := filepath.Join(dbDir, ".tmp_write_test")
		file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("database directory is not writable: %v", err)
		}
		file.Close()
		os.Remove(tempFile)
	}

	// Open database connection
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection and write permissions by performing a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction (database might be locked or not writable): %v", err)
	}
	tx.Rollback() // Always rollback this test transaction

	// Create students table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		linkedin_profile TEXT,
		phone TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Printf("Successfully connected to database: %s", dbPath)
	if fileExists {
		log.Printf("Using existing database file")
	} else {
		log.Printf("Created new database file")
	}

	return nil
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO students (name, email, linkedin_profile, phone) VALUES (?, ?, ?, ?)",
		student.Name, student.Email, student.LinkedinProfile, student.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	student.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student

	err := db.QueryRow("SELECT id, name, email, linkedin_profile, phone FROM students WHERE id = ?",
		params["id"]).Scan(&student.ID, &student.Name, &student.Email, &student.LinkedinProfile, &student.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, linkedin_profile, phone FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.LinkedinProfile, &student.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE students SET name = ?, email = ?, linkedin_profile = ?, phone = ? WHERE id = ?",
		student.Name, student.Email, student.LinkedinProfile, student.Phone, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	student.ID = parseInt(params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("DELETE FROM students WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseInt(s string) int {
	var i int
	json.Unmarshal([]byte(s), &i)
	return i
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(swaggerContent)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	response := WelcomeResponse{
		Message:     "Welcome to Students API",
		Description: "A RESTful API for managing student records with CRUD operations",
		Swagger:     "/swagger",
		Version:     "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to run the server on")
	dbPath := flag.String("db", "./students.db", "Path to SQLite database file")
	flag.Parse()

	// Initialize database
	if err := initDB(*dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	// Welcome endpoint
	router.HandleFunc("/", welcomeHandler).Methods("GET")

	// CRUD endpoints
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students", getAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	// Swagger endpoint
	router.HandleFunc("/swagger", serveSwagger).Methods("GET")

	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("Using database: %s", *dbPath)
	log.Printf("Starting server on port %d...", *port)
	log.Printf("To use a different port, restart with: %s --port <port-number>", os.Args[0])
	
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		if err.Error() == fmt.Sprintf("listen tcp %s: bind: address already in use", serverAddr) {
			log.Printf("Port %d is already in use. Try using a different port with: %s --port <port-number>", *port, os.Args[0])
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

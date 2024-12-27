# Student Management API

A RESTful API built with Go and SQLite that provides CRUD operations for managing student records. The API stores basic student information including name, email, LinkedIn profile, and phone number.

## Project Structure

- `main.go` - Main API server implementation
- `init_db.go` - Database initialization script with sample data
- `students.db` - SQLite database file (created when running the application)

## Prerequisites

- Go 1.16 or higher
- SQLite3

## Getting Started

1. Clone the repository:
``` 
git clone <your-repo-url>
cd <repo-directory>
```

2. Install dependencies:
``` 
go mod tidy
```

3. Initialize the database with sample data:
``` 
go run init_db.go
```

4. Start the API server:
``` 
go run main.go
```

The server will start on http://localhost:8080

## API Endpoints and Usage

### Get All Students
``` 
curl http://localhost:8080/students
```

### Get a Specific Student
``` 
curl http://localhost:8080/students/1
```

### Create a New Student
``` 
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Brown",
    "email": "alice.brown@example.com",
    "linkedin_profile": "https://linkedin.com/in/alicebrown",
    "phone": "+1-555-987-6543"
  }'
```

### Update a Student
``` 
curl -X PUT http://localhost:8080/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.doe@example.com",
    "linkedin_profile": "https://linkedin.com/in/johndoe",
    "phone": "+1-555-123-4567"
  }'
```

### Delete a Student
``` 
curl -X DELETE http://localhost:8080/students/1
```

## API Response Codes

- 200: Success
- 201: Created (For successful POST requests)
- 204: No Content (For successful DELETE requests)
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

## Data Structure

Student record fields:
- id (integer, auto-generated)
- name (string, required)
- email (string, required, unique)
- linkedin_profile (string, optional)
- phone (string, optional)

## Sample Data

The initialization script (`init_db.go`) populates the database with three sample students:
1. John Doe
2. Jane Smith
3. Bob Johnson

You can modify the sample data by editing the `init_db.go` file.

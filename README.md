# Student Management API

A RESTful API built with Go and SQLite that provides CRUD operations for managing student records. The API stores basic student information including name, email, LinkedIn profile, and phone number.

## Project Structure

- `main.go` - Main API server implementation
- `init_db.go` - Database initialization script with sample data
- `swagger.yaml` - OpenAPI/Swagger specification for the API
- `students.db` - Default SQLite database file (created when running the application)

## Prerequisites

- Go 1.16 or higher
- SQLite3

## Getting Started

1. Clone the repository:
``` 
git clone https://github.com/yourusername/students-api-golang.git
cd students-api-golang
```

2. Install dependencies:
``` 
go mod tidy
```

3. Initialize the database with sample data:
``` 
# Use default database location (./students.db)
go run init_db.go

# Or specify a custom database location
go run init_db.go --db /path/to/your/students.db
```

4. Start the API server:
``` 
# Use default port (8080) and default database
go run main.go

# Use custom port and/or custom database
go run main.go --port 8082 --db /path/to/your/students.db
```

The server will start on the specified port (default: 8080) using the specified database (default: ./students.db). 
If the default port is already in use, you'll receive a helpful message suggesting to use a different port with the `--port` flag.
The database file and its directory will be created automatically if they don't exist.

## API Documentation

The API documentation is available in OpenAPI/Swagger format. You can access it at:
# If using default port
curl http://localhost:8080/swagger

# If using custom port (e.g., 8082)
curl http://localhost:8082/swagger

You can also use this URL with Swagger UI or other OpenAPI tools to visualize and interact with the API documentation.

## API Endpoints and Usage

Replace `<port>` in the examples below with your chosen port number (default: 8080)

### Root Endpoint
``` 
curl http://localhost:<port>/
```

Response:
``` 
{
    "message": "Welcome to Students API",
    "description": "A RESTful API for managing student records with CRUD operations",
    "swagger_url": "/swagger",
    "version": "1.0.0"
}
```

### Get All Students
``` 
curl http://localhost:<port>/students
```

### Get a Specific Student
``` 
curl http://localhost:<port>/students/1
```

### Create a New Student
``` 
curl -X POST http://localhost:<port>/students \
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
curl -X PUT http://localhost:<port>/students/1 \
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
curl -X DELETE http://localhost:<port>/students/1
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

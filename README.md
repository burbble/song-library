# Song Library

A REST API service for managing a song library. Allows storing and managing information about songs, including band name, song title, lyrics, release date and links.

## Tech Stack

- Go 1.22.3
- PostgreSQL 15
- Docker & Docker Compose
- Gin Web Framework
- Swagger for API documentation
- Zap Logger
- Make for automation

## Prerequisites

- Docker and Docker Compose
- Make
- Go 1.22.3


1. Create .env file from example:

2. Start the project:

```bash
make start
```

The service will be available at: http://localhost:8080

Swagger UI: http://localhost:8080/swagger/index.html
PgAdmin: http://localhost:5050 (login: admin@admin.com, password: admin)

## Make Commands Reference

### Core Commands
- `make start` - Full application startup: runs migrations, seeds database, builds and starts all services
- `make up` - Start all services in docker containers
- `make down` - Stop and remove all containers and volumes
- `make build` - Build docker images
- `make up --build` - Rebuild and start all services (useful after code changes)

### Database Commands
- `make postgres` - Start PostgreSQL container and wait for it to be ready
- `make recreate-db` - Drop and recreate the database (requires postgres to be running)
- `make migrate` - Recreate database and apply all migrations
- `make seed` - Populate database with initial test data
- `make reset-db` - Full database reset: recreate, migrate and seed

### Service Management
- `make restart-app` - Restart only the application container (useful during development)

### Logging Commands
View logs using the `logs` command with optional service parameter:
- `make logs` - View logs from all services
- `make logs service=app` - View only application logs
- `make logs service=postgres` - View only database logs

### Command Details

#### Database Commands

## Project Structure

```
song-library/
├── cmd/                    # Application entry point
├── docs/                   # Swagger documentation
├── internal/               # Internal application code
│   ├── application/       # Business logic and DTOs
│   ├── config/           # Application configuration
│   ├── domain/           # Domain models
│   ├── infrastructure/   # Repository implementations
│   └── interfaces/       # HTTP handlers
├── migrations/            # SQL migrations
├── pkg/                   # Common packages
├── scripts/              # DB scripts
└── docker-compose.yml    # Docker configuration
```

## API Endpoints

### Songs

- `GET /api/v1/songs` - Get list of songs with filtering and pagination
- `POST /api/v1/songs` - Create new song
- `GET /api/v1/songs/{id}` - Get song by ID
- `PUT /api/v1/songs/{id}` - Update song
- `DELETE /api/v1/songs/{id}` - Delete song
- `GET /api/v1/songs/{id}/text` - Get song text with verse pagination

## Development

### Local Setup

1. Start PostgreSQL:
```bash
make postgres
```

2. Apply migrations:
```bash
make migrate
```

3. Seed test data:
```bash
make seed
```

### Configuration

Main settings are in `.env` file:

- `SERVER_PORT` - Server port (default 8080)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `LOG_LEVEL` - Logging level

## Logging

The application uses Zap Logger for structured logging. Log levels:
- debug - Detailed debugging information
- info - Informational messages
- warn - Warning messages
- error - Error messages

## Database

PostgreSQL is used as the main data store. Database schema:

```sql
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## API Documentation

Full API documentation is available via Swagger UI at http://localhost:8080/swagger/index.html when the service is running.

### Example Requests

#### Create a new song
```bash
curl -X POST http://localhost:8080/api/v1/songs \
  -H "Content-Type: application/json" \
  -d '{
    "group": "Queen",
    "song": "Bohemian Rhapsody"
  }'
```

#### Get songs list with pagination
```bash
curl "http://localhost:8080/api/v1/songs?page=1&page_size=10"
```

## Error Handling

The API uses standard HTTP status codes and returns errors in the following format:

```json
{
    "error": "Error message description"
}
```

Common status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

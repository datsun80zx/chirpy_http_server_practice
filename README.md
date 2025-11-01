# Chirpy HTTP Server

A Twitter-like social media platform built with Go, featuring user authentication, posts (chirps), and a premium subscription system.

## Features

- 🔐 User authentication with JWT tokens
- 🔄 Refresh token system for extended sessions
- 📝 Create, read, and delete chirps (posts)
- ✨ Profanity filter for chirps
- 👑 Premium user upgrades (Chirpy Red)
- 🔒 Secure password hashing with Argon2
- 🗄️ PostgreSQL database with sqlc-generated type-safe queries
- 🪝 Webhook support for payment processing

## Tech Stack

- **Language**: Go 1.23.4
- **Database**: PostgreSQL with [sqlc](https://sqlc.dev/)
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: Argon2id
- **Migration Tool**: Goose
- **Dependencies**:
  - `github.com/google/uuid` - UUID generation
  - `github.com/lib/pq` - PostgreSQL driver
  - `github.com/joho/godotenv` - Environment variable management
  - `github.com/alexedwards/argon2id` - Password hashing
  - `github.com/golang-jwt/jwt/v5` - JWT tokens

## Prerequisites

- Go 1.23.4 or higher
- PostgreSQL database
- [sqlc](https://sqlc.dev/) (for regenerating database code)
- [goose](https://github.com/pressly/goose) (for database migrations)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd chirpy_http_server_practice
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
Create a `.env` file in the root directory:
```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
TOKEN_STRING=your-jwt-secret-key
POLKA_KEY=your-polka-api-key
```

4. Run database migrations:
```bash
goose -dir sql/schema postgres "your-db-url" up
```

5. Generate sqlc code (if schema changes):
```bash
sqlc generate
```

6. Run the server:
```bash
go run .
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080
```

---

### Authentication

Most endpoints require authentication via JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

---

## API Endpoints

### Health Check

#### `GET /api/healthz`
Check if the API is running.

**Response**: `200 OK`
```
OK
```

---

### User Management

#### `POST /api/users`
Create a new user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2025-11-01T10:00:00Z",
  "updated_at": "2025-11-01T10:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

---

#### `PUT /api/users`
Update user email and password.

**Authentication**: Required (Bearer token)

**Request Body**:
```json
{
  "email": "newemail@example.com",
  "password": "newsecurepassword123"
}
```

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2025-11-01T10:00:00Z",
  "updated_at": "2025-11-01T10:15:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

---

### Authentication Endpoints

#### `POST /api/login`
Authenticate user and receive access token and refresh token.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2025-11-01T10:00:00Z",
  "updated_at": "2025-11-01T10:00:00Z",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6...",
  "is_chirpy_red": false
}
```

**Notes**:
- Access token expires in 1 hour
- Refresh token expires in 60 days

---

#### `POST /api/refresh`
Refresh an expired access token using a refresh token.

**Authentication**: Required (Bearer token with refresh token)

**Headers**:
```
Authorization: Bearer <your-refresh-token>
```

**Response**: `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

#### `POST /api/revoke`
Revoke a refresh token (logout).

**Authentication**: Required (Bearer token with refresh token)

**Headers**:
```
Authorization: Bearer <your-refresh-token>
```

**Response**: `204 No Content`

---

### Chirps (Posts)

#### `POST /api/chirps`
Create a new chirp.

**Authentication**: Required (Bearer token)

**Request Body**:
```json
{
  "body": "This is my first chirp!"
}
```

**Response**: `201 Created`
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "created_at": "2025-11-01T10:30:00Z",
  "updated_at": "2025-11-01T10:30:00Z",
  "body": "This is my first chirp!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Validation**:
- Maximum 140 characters
- Profanity filter replaces offensive words with `****`
  - Filtered words: "kerfuffle", "sharbert", "fornax" (case-insensitive)

---

#### `GET /api/chirps`
Get all chirps with optional filtering and sorting.

**Query Parameters**:
- `author_id` (optional): Filter chirps by user ID
- `sort` (optional): Sort order
  - `asc` (default): Oldest first
  - `desc`: Newest first

**Examples**:
```
GET /api/chirps
GET /api/chirps?sort=desc
GET /api/chirps?author_id=123e4567-e89b-12d3-a456-426614174000
GET /api/chirps?author_id=123e4567-e89b-12d3-a456-426614174000&sort=desc
```

**Response**: `200 OK`
```json
[
  {
    "id": "456e7890-e89b-12d3-a456-426614174001",
    "created_at": "2025-11-01T10:30:00Z",
    "updated_at": "2025-11-01T10:30:00Z",
    "body": "This is my first chirp!",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  },
  {
    "id": "789e0123-e89b-12d3-a456-426614174002",
    "created_at": "2025-11-01T10:35:00Z",
    "updated_at": "2025-11-01T10:35:00Z",
    "body": "Another chirp here!",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
]
```

---

#### `GET /api/chirps/{chirpID}`
Get a specific chirp by ID.

**Response**: `200 OK`
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "created_at": "2025-11-01T10:30:00Z",
  "updated_at": "2025-11-01T10:30:00Z",
  "body": "This is my first chirp!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Error Response**: `404 Not Found`
```json
{
  "error": "Chirp not found"
}
```

---

#### `DELETE /api/chirps/{chirpID}`
Delete a chirp (only the author can delete their own chirps).

**Authentication**: Required (Bearer token)

**Response**: `204 No Content`

**Error Responses**:
- `403 Forbidden`: User is not the author of the chirp
- `404 Not Found`: Chirp doesn't exist

---

### Premium Features

#### `POST /api/polka/webhooks`
Webhook endpoint for Polka payment processor to upgrade users to Chirpy Red.

**Authentication**: API Key required

**Headers**:
```
Authorization: ApiKey <your-polka-api-key>
```

**Request Body**:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

**Response**: `204 No Content`

**Notes**:
- Only processes `user.upgraded` events
- Other event types return `204 No Content` without action
- User must exist in the database

---

### Admin Endpoints

#### `GET /admin/metrics`
View the number of times the file server has been accessed.

**Response**: `200 OK`
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

---

#### `POST /admin/reset`
Reset the database and hit counter (development only).

**Response**: `200 OK`
```
Hits reset to 0
```

**Error Response**: `403 Forbidden` (if not in development environment)

**Notes**:
- Only available when `PLATFORM=dev`
- Deletes all users (and cascades to chirps and refresh tokens)

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

Common status codes:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource doesn't exist
- `500 Internal Server Error`: Server-side error

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    is_chirpy_red BOOLEAN DEFAULT false
);
```

### Chirps Table
```sql
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Test specific package:
```bash
go test ./internal/auth
```

---

## Project Structure

```
.
├── main.go                      # Application entry point
├── validate.go                  # Legacy validation code (commented out)
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── sqlc.yaml                    # sqlc configuration
├── index.html                   # Static homepage
├── .gitignore                   # Git ignore rules
├── README.md                    # This file
│
├── sql/
│   ├── schema/                  # Database migrations
│   │   ├── 001_users.sql
│   │   ├── 002_chirps.sql
│   │   ├── 003_users.sql
│   │   ├── 004_refresh_tokens.sql
│   │   └── 005_users.sql
│   └── queries/                 # SQL queries for sqlc
│       ├── users.sql
│       ├── chirps.sql
│       └── refresh_tokens.sql
│
└── internal/
    ├── config.go                # Application configuration
    │
    ├── auth/                    # Authentication package
    │   ├── jwt.go               # JWT creation and validation
    │   ├── jwt_test.go          # JWT tests
    │   ├── hasher.go            # Password hashing
    │   ├── hasher_test.go       # Hashing tests
    │   └── api_key.go           # API key extraction
    │
    ├── api/                     # API handlers
    │   ├── api.go               # Handler initialization
    │   ├── json.go              # JSON response helpers
    │   ├── readiness.go         # Health check handler
    │   ├── users.go             # User management endpoints
    │   ├── login.go             # Login endpoint
    │   ├── refresh.go           # Token refresh endpoint
    │   ├── revoke.go            # Token revoke endpoint
    │   ├── chirps.go            # Chirp endpoints
    │   └── reset.go             # Admin reset endpoint
    │
    └── database/                # Generated database code (sqlc)
        ├── db.go                # Database interface
        ├── models.go            # Type definitions
        ├── users.sql.go         # User queries
        ├── chirps.sql.go        # Chirp queries
        └── refresh_tokens.sql.go # Refresh token queries
```

---

## Security Considerations

- ✅ Passwords are hashed using Argon2id (industry-standard)
- ✅ JWT tokens expire after 1 hour
- ✅ Refresh tokens expire after 60 days
- ✅ Refresh tokens can be revoked
- ✅ SQL injection protection via prepared statements (sqlc)
- ✅ HTTPS recommended for production
- ⚠️ Store `TOKEN_STRING` and `POLKA_KEY` securely
- ⚠️ Never commit `.env` file to version control

---

## Development

### Adding New Database Tables

1. Create migration file in `sql/schema/`:
```bash
goose -dir sql/schema create table_name sql
```

2. Write up and down migrations

3. Run migration:
```bash
goose -dir sql/schema postgres "your-db-url" up
```

4. Create queries in `sql/queries/`

5. Regenerate sqlc code:
```bash
sqlc generate
```

### Code Style

- Follow standard Go formatting: `go fmt ./...`
- Run linter: `go vet ./...`
- Write tests for new features

---

## License

This project is for educational purposes.

---

## Contributing

This is a practice project, but suggestions are welcome!

---

## Support

For issues or questions, please open an issue on the repository.
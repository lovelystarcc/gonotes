# GoNotes

Minimal REST API for note management with JWT authentication, written in Go with **chi** and **render**. Data is stored in a **SQLite** database.

---

## Features

- **User Authentication:** JWT-based registration and login
- **Note Management:**
  - Create note: `POST /notes`
  - Get note by ID: `GET /notes/{id}`
  - Delete note by ID: `DELETE /notes/{id}`
  - List all notes: `GET /notes`
- **Security:** Password hashing with bcrypt, JWT token validation
- **CORS Support:** Configured for local development
- **Structured Logging:** Colorful console output with slog
- **Persistent Storage:** Users and notes stored in SQLite

---

## Requirements

- **Go:** 1.20+
- **SQLite:** via [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- **Migrate CLI:** for database migrations

---

## Getting Started

### Configuration
Create a config file at `configs/local.yaml`:

```yaml
env: "local"
secret_key: "your-secret-key-here"
storage_path: "storage/storage.db"
http_server:
  address: ":8080"
  timeout: 4s
  idle_timeout: 60s
  user: "admin"
  password: "password"
```

### Database Migrations
```bash
# Create new migration
make migrate-create name=init_schema

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Using Makefile
- **Run server:**
  ```bash
  make run
  ```
- **Build binary:**
  ```bash
  make build
  ```
- **Run tests:**
  ```bash
  make test
  ```
- **Format code:**
  ```bash
  make fmt
  ```
- **Run linter:**
  ```bash
  make lint
  ```
- **Clean build artifacts:**
  ```bash
  make clean
  ```

---

## API

### Authentication Endpoints

#### Register
- **URL:** `POST /auth/register`
- **Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
- **Response 201:**
```json
{
  "id": 1,
  "email": "user@example.com"
}
```

#### Login
- **URL:** `POST /auth/login`
- **Body:** Same as register
- **Response 200:**
```json
{
  "email": "user@example.com",
  "token": "jwt-token-here"
}
```

### Note Endpoints (All require authentication)

#### Create Note
- **URL:** `POST /notes`
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Body:**
```json
{
  "title": "My Note",
  "content": "Note content here"
}
```
- **Response 201:**
```json
{
  "id": 1,
  "title": "My Note",
  "content": "Note content here",
  "created_at": "2023-01-01T00:00:00Z"
}
```

#### Get Note
- **URL:** `GET /notes/{id}`
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Response 200:**
```json
{
  "id": 1,
  "title": "My Note",
  "content": "Note content here",
  "created_at": "2023-01-01T00:00:00Z"
}
```

#### Delete Note
- **URL:** `DELETE /notes/{id}`
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Response 200:** Returns deleted note

#### List Notes
- **URL:** `GET /notes`
- **Headers:** `Authorization: Bearer <jwt-token>`
- **Response 200:**
```json
[
  {
    "id": 1,
    "title": "My Note",
    "content": "Note content here",
    "created_at": "2023-01-01T00:00:00Z"
  }
]
```

---

## Usage Examples

### 1. Register a new user
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 2. Login and get token
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 3. Create a note (with JWT token)
```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"First Note","content":"This is my first note"}'
```

### 4. Get all notes
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/notes
```

---

## Error Model

- **Format:**
```json
{
  "message": "error description"
}
```
- **Common Status Codes:**
  - **400 Bad Request:** Invalid input data
  - **401 Unauthorized:** Missing or invalid JWT token
  - **404 Not Found:** Note not found
  - **409 Conflict:** User already exists
  - **500 Internal Server Error:** Server error

---

## Implementation Details

- **Router:** [github.com/go-chi/chi](https://github.com/go-chi/chi)
- **Rendering:** [github.com/go-chi/render](https://github.com/go-chi/render)
- **Authentication:** JWT with [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- **Password Hashing:** bcrypt via golang.org/x/crypto
- **Storage:** SQLite with repository pattern
- **Logging:** Structured logging with slog and custom pretty handler
- **Configuration:** YAML config with cleanenv

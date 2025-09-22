# GoNotes

Minimal REST API for note management, written in Go with chi and render. Data is stored in memory and resets on server restart.

---

## Features

- **Create note:** POST /notes
- **Get note by ID:** GET /notes/{id}
- **Delete note by ID:** DELETE /notes/{id}
- **List all notes:** GET /notes
- **JSON I/O with basic validation:** text is required

---

## Requirements

- **Go:** 1.20+

---

## Getting started

### Using Makefile
- **Run server:**
  ```bash
  make run
  ```
- **Build binary:**
  ```bash
  make build
  ```
  - Output binary name is controlled by:
    ```makefile
    BINARY_NAME=app
    ```
  - Result: ./app
- **Run tests:**
  ```bash
  make test
  ```
- **Clean build artifacts:**
  ```bash
  make clean
  ```

### Without Make
- **Install deps:**
  ```bash
  go mod tidy
  ```
- **Run:**
  ```bash
  go run cmd/gonotes/main.go
  ```
- **Build:**
  ```bash
  go build -o app main.go
  ```

- **Server URL:** http://localhost:8080

---

## API

### Endpoints
| Method | Path        | Description         |
|-------:|-------------|---------------------|
| POST   | /notes      | Create a note       |
| GET    | /notes/{id} | Get a note by ID    |
| DELETE | /notes/{id} | Delete a note by ID |
| GET    | /notes      | List all notes      |

> ID is an integer and auto-incremented in memory.

### Create a note
- **Request body:**
```json
{ "text": "First note" }
```
- **Response 201:**
```json
{ "id": 1, "text": "First note" }
```
- **Validation:** missing text → 400

Example:
```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"text":"First note"}'
```

### Get a note
- **Response 200:**
```json
{ "id": 1, "text": "First note" }
```
Example:
```bash
curl http://localhost:8080/notes/1
```

### Delete a note
- **Response:** 204 No Content

Example:
```bash
curl -X DELETE http://localhost:8080/notes/1
```

### List notes
- **Response 200:**
```json
[
  { "id": 1, "text": "First note" }
]
```
Example:
```bash
curl http://localhost:8080/notes
```

---

## Error model

- **Format:**
```json
{ "message": "error description" }
```
- **Examples:**
  - **400 Bad Request:** invalid ID format, missing text
  - **404 Not Found:** note not found

---

## Implementation notes

- **Router:** github.com/go-chi/chi
- **Rendering:** github.com/go-chi/render
- **Storage:** in-memory map (non-persistent)
- **Port:** 8080 (hard-coded in main.go)
- **Concurrency:** demo-only; not designed for concurrent writes

---

## License

MIT. See LICENSE.

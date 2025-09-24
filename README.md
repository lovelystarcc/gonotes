# GoNotes

Minimal REST API for note management, written in Go with **chi** and **render**.
Data is stored in a **SQLite** database.

---

## Features

- **Create note:** `POST /notes`
- **Get note by ID:** `GET /notes/{id}`
- **Delete note by ID:** `DELETE /notes/{id}`
- **List all notes:** `GET /notes`
- **JSON I/O with basic validation:** `text` is required
- **Persistent storage:** notes are saved in `storage/storage.db`

---

## Requirements

- **Go:** 1.20+
- **SQLite:** via [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

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
  - Result: `./app`

- **Clean build artifacts:**
  ```bash
  make clean
  ```

## API

### Endpoints
| Method | Path        | Description         |
|-------:|-------------|---------------------|
| POST   | /notes      | Create a note       |
| GET    | /notes/{id} | Get a note by ID    |
| GET    | /notes      | List all notes      |
| DELETE | /notes/{id} | Delete a note by ID |

> ID is an integer and auto-incremented by SQLite.

---

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

---

### Get a note
- **Response 200:**
```json
{ "id": 1, "text": "First note" }
```
Example:
```bash
curl http://localhost:8080/notes/1
```

---

### Delete a note
- **Response:** 204 No Content (no body)
  *or* 200 OK with deleted note if API is configured that way.

Example:
```bash
curl -X DELETE http://localhost:8080/notes/1
```

---

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

- **Router:** [github.com/go-chi/chi](https://github.com/go-chi/chi)
- **Rendering:** [github.com/go-chi/render](https://github.com/go-chi/render)
- **Storage:** SQLite via repository interface (`internal/storage.NoteRepository`)

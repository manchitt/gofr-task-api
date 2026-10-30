# GoFr Task API — Sample Observable Microservice

A production-ready sample REST microservice built with [GoFr](https://gofr.dev) demonstrating:

- **Built-in SQLite support & Datasource management**
- **Automated versioned schema migrations** (`migration.Migrate`)
- **Full CRUD REST endpoints** (`POST`, `GET`, `PUT`, `DELETE /tasks`)
- **Built-in observability**: OpenTelemetry distributed tracing, structured logging with correlation IDs, Prometheus metrics
- **Custom application metrics**: `tasks_created_total`
- **Native health & liveness checks**: `/.well-known/health` and `/.well-known/alive`
- **End-to-end integration tests** with GoFr's `testutil`

---

## Prerequisites

- Go 1.26+

---

## Configuration (`configs/.env`)

```ini
APP_NAME=task-api
HTTP_PORT=9000
METRICS_PORT=2121
DB_DIALECT=sqlite
DB_NAME=tasks.db
```

---

## Running the Application

```bash
go run main.go
```

The service will:
1. Load configuration from `configs/.env`.
2. Connect to SQLite database `tasks.db`.
3. Automatically execute pending migrations in `migrations/`.
4. Start the REST API on port `9000`.
5. Start the Prometheus metrics server on port `2121`.

---

## API Endpoints

| Method | Path | Description | Status Code |
|---|---|---|---|
| `POST` | `/tasks` | Create a new task | `201 Created` |
| `GET` | `/tasks` | List all tasks (`?status=pending`) | `200 OK` |
| `GET` | `/tasks/{id}` | Get task by ID | `200 OK` / `404 Not Found` |
| `PUT` | `/tasks/{id}` | Update task | `200 OK` / `404 Not Found` |
| `DELETE` | `/tasks/{id}` | Delete task | `204 No Content` / `404 Not Found` |
| `GET` | `/.well-known/health` | Service & database health check | `200 OK` |
| `GET` | `/.well-known/alive` | Service liveness probe | `200 OK` |

---

## Running Integration Tests

```bash
go test -v ./...
```

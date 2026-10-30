package main

import (
	"database/sql"
	"errors"
	"fmt"

	"gofr-task-api/migrations"
	"gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func main() {
	app := gofr.New()

	// Register database migrations
	app.Migrate(migrations.All())

	// Register Task REST endpoints
	app.POST("/tasks", createTask)
	app.GET("/tasks", listTasks)
	app.GET("/tasks/{id}", getTask)
	app.PUT("/tasks/{id}", updateTask)
	app.DELETE("/tasks/{id}", deleteTask)

	app.Run()
}

func createTask(c *gofr.Context) (any, error) {
	var task Task
	if err := c.Bind(&task); err != nil {
		return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"body"}}
	}

	if task.Title == "" {
		return nil, gofrHTTP.ErrorMissingParam{Params: []string{"title"}}
	}

	if task.Status == "" {
		task.Status = "pending"
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	query := `INSERT INTO tasks (title, description, status, priority) VALUES (?, ?, ?, ?)`
	res, err := c.SQL.ExecContext(c, query, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	task.ID = id

	// Increment custom Prometheus metric
	c.Metrics().IncrementCounter(c, "tasks_created_total")

	return task, nil
}

func listTasks(c *gofr.Context) (any, error) {
	statusFilter := c.Param("status")

	var (
		rows *sql.Rows
		err  error
	)

	if statusFilter != "" {
		query := `SELECT id, title, description, status, priority, created_at FROM tasks WHERE status = ? ORDER BY id DESC`
		rows, err = c.SQL.QueryContext(c, query, statusFilter)
	} else {
		query := `SELECT id, title, description, status, priority, created_at FROM tasks ORDER BY id DESC`
		rows, err = c.SQL.QueryContext(c, query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func getTask(c *gofr.Context) (any, error) {
	id := c.PathParam("id")
	if id == "" {
		return nil, gofrHTTP.ErrorMissingParam{Params: []string{"id"}}
	}

	var t Task
	query := `SELECT id, title, description, status, priority, created_at FROM tasks WHERE id = ?`
	err := c.SQL.QueryRowContext(c, query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, gofrHTTP.ErrorEntityNotFound{Name: "task", Value: id}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}

	return t, nil
}

func updateTask(c *gofr.Context) (any, error) {
	id := c.PathParam("id")
	if id == "" {
		return nil, gofrHTTP.ErrorMissingParam{Params: []string{"id"}}
	}

	var update Task
	if err := c.Bind(&update); err != nil {
		return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"body"}}
	}

	query := `UPDATE tasks SET title = COALESCE(NULLIF(?, ''), title),
	                           description = COALESCE(NULLIF(?, ''), description),
	                           status = COALESCE(NULLIF(?, ''), status),
	                           priority = COALESCE(NULLIF(?, ''), priority)
	          WHERE id = ?`

	res, err := c.SQL.ExecContext(c, query, update.Title, update.Description, update.Status, update.Priority, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, gofrHTTP.ErrorEntityNotFound{Name: "task", Value: id}
	}

	return getTask(c)
}

func deleteTask(c *gofr.Context) (any, error) {
	id := c.PathParam("id")
	if id == "" {
		return nil, gofrHTTP.ErrorMissingParam{Params: []string{"id"}}
	}

	query := `DELETE FROM tasks WHERE id = ?`
	res, err := c.SQL.ExecContext(c, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, gofrHTTP.ErrorEntityNotFound{Name: "task", Value: id}
	}

	return map[string]string{"message": fmt.Sprintf("task %s deleted successfully", id)}, nil
}

package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTasksTable = `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    priority TEXT NOT NULL DEFAULT 'medium',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

func createTasks() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTasksTable)
			return err
		},
	}
}

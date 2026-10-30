package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

// All returns all migrations registered in chronological order.
func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		1730300000: createTasks(),
	}
}

package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedBoards(db *sql.DB) error {

	var projKanbanID string

	if err := db.QueryRow("SELECT id FROM projects WHERE name = $1", "Kanban Project").Scan(&projKanbanID); err != nil {
		return fmt.Errorf("kanban project not found: %v", err)
	}

	now := time.Now()
	_, err := db.Exec(`
		insert into boards (id, project_id, name, created_at, updated_at) values
		($1, $2, 'Backlog', $3, $4),
		($5, $6, 'In Progress', $7, $8),
		($9, $10, 'Done', $11, $12)
	`,
		uuid.NewString(), projKanbanID, now, now,
		uuid.NewString(), projKanbanID, now, now,
		uuid.NewString(), projKanbanID, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding boards: %v", err)
	}

	fmt.Println("✅ Boards seeder successfully")
	return nil
}

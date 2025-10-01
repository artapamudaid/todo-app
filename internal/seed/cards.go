package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedCards(db *sql.DB) error {
	var boardBacklog, boardInProgress, boardDone string

	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "Backlog").Scan(&boardBacklog); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}
	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "In Progress").Scan(&boardInProgress); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}
	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "Done").Scan(&boardDone); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}

	now := time.Now()
	_, err := db.Exec(`
		insert into cards (id, board_id, name, is_closed, created_at, updated_at) values
		($1, $2, 'Setup project', false, $3, $4),
		($5, $6, 'Build API', false, $7, $8),
		($9, $10, 'Design DB schema', true, $11, $12)
	`,
		uuid.NewString(), boardBacklog, now, now,
		uuid.NewString(), boardInProgress, now, now,
		uuid.NewString(), boardDone, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding cards: %v", err)
	}

	fmt.Println("✅ Cards seeder successfully")
	return nil
}

package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedCards(db *sql.DB) error {
	var boardBacklog, boardInProgress, boardDone, userId string

	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "Backlog").Scan(&boardBacklog); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}
	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "In Progress").Scan(&boardInProgress); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}
	if err := db.QueryRow("SELECT id FROM boards WHERE name = $1", "Done").Scan(&boardDone); err != nil {
		return fmt.Errorf("kanban project -> board backlog not found: %v", err)
	}

	if err := db.QueryRow("SELECT id FROM users WHERE name = $1", "Admin").Scan(&userId); err != nil {
		return fmt.Errorf("user Admin -> user.name not found: %v", err)
	}

	now := time.Now()
	_, err := db.Exec(`
		insert into cards (id, board_id, name, is_closed, created_at, updated_at, user_id) values
		($1, $2, 'Setup project', false, $3, $4, $5),
		($6, $7, 'Build API', false, $8, $9, $10),
		($11, $12, 'Design DB schema', false, $13, $14, $15)
	`,
		uuid.NewString(), boardBacklog, now, now, userId,
		uuid.NewString(), boardInProgress, now, now, userId,
		uuid.NewString(), boardDone, now, now, userId,
	)

	if err != nil {
		return fmt.Errorf("failed seeding cards: %v", err)
	}

	fmt.Println("✅ Cards seeder successfully")
	return nil
}

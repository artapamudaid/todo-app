package seed

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func SeedProjectUsers(db *sql.DB) error {
	var adminID, projKanbanID, projHrID string

	// Ambil ID user by email
	if err := db.QueryRow("SELECT id FROM users WHERE email = $1", "admin@todo.app").Scan(&adminID); err != nil {
		return fmt.Errorf("admin user not found: %v", err)
	}

	// Ambil ID project by name
	if err := db.QueryRow("SELECT id FROM projects WHERE name = $1", "Kanban Project").Scan(&projKanbanID); err != nil {
		return fmt.Errorf("kanban project not found: %v", err)
	}
	if err := db.QueryRow("SELECT id FROM projects WHERE name = $1", "HR System").Scan(&projHrID); err != nil {
		return fmt.Errorf("hr project not found: %v", err)
	}

	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO project_users (id, project_id, user_id, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5),
		($6, $7, $8, $9, $10)
	`,
		uuid.NewString(), projKanbanID, adminID, now, now,
		uuid.NewString(), projHrID, adminID, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding project users: %v", err)
	}

	fmt.Println("✅ Project Users seeder successfully")
	return nil
}

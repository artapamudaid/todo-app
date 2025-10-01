package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedProjects(db *sql.DB) error {
	now := time.Now()
	_, err := db.Exec(`
		insert into projects (id, name, created_at, updated_at) values
		($1, 'Kanban Project', $2, $3),
		($4, 'HR System', $5, $6)
	`, uuid.NewString(), now, now,
		uuid.NewString(), now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding projects: %v", err)
	}

	fmt.Println("✅ Projects seeder successfully")
	return nil
}

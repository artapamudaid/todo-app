package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedProjects(db *sql.DB) error {
	var departmentId string

	if err := db.QueryRow("SELECT id FROM departments WHERE name = $1", "IT").Scan(&departmentId); err != nil {
		return fmt.Errorf("department IT -> departments.name not found: %v", err)
	}

	now := time.Now()
	_, err := db.Exec(`
		insert into projects (id, name, created_at, updated_at, department_id) values
		($1, 'Kanban Project', $2, $3, $4),
		($5, 'HR System', $6, $7, $8)
	`, uuid.NewString(), now, now, departmentId,
		uuid.NewString(), now, now, departmentId,
	)

	if err != nil {
		return fmt.Errorf("failed seeding projects: %v", err)
	}

	fmt.Println("✅ Projects seeder successfully")
	return nil
}

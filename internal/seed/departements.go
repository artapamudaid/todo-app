package seed

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

func SeedDepartments(db *sql.DB) error {
	now := time.Now()
	_, err := db.Exec(`
		insert into departments (id, name, created_at, updated_at) values
		($1, 'IT', $2, $3),
		($4, 'HR', $5, $6)
	`,
		uuid.NewString(), now, now,
		uuid.NewString(), now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding departements: %v", err)
	}

	fmt.Println("✅ Departements seeder successfully")
	return nil
}

package seed

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func SeedRoles(db *sql.DB) error {
	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO roles (id, name, created_at, updated_at)
		VALUES
			($1, 'Admin', $2, $3),
			($4, 'User',  $5, $6)
		ON CONFLICT (id) DO NOTHING
	`,
		uuid.NewString(), now, now,
		uuid.NewString(), now, now,
	)

	if err != nil {
		return fmt.Errorf("failed seeding roles: %v", err)
	}

	fmt.Println("✅ Roles seeder successfully")
	return nil
}

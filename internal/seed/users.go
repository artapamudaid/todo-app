package seed

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app/internal/util/helper"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func SeedUsers(db *sql.DB) error {
	// Ambil role_id & department_id
	var roleID string
	var deptID string

	err := db.QueryRow("SELECT id FROM roles WHERE name = $1 LIMIT 1", "Admin").Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	err = db.QueryRow("SELECT id FROM departments WHERE name = $1 LIMIT 1", "IT").Scan(&deptID)
	if err != nil {
		return fmt.Errorf("department not found: %v", err)
	}

	// Hash password
	hashedPassword, _ := helper.HashPassword("password123")

	// Data user
	id := uuid.NewString()
	name := "Super Admin"
	email := "admin@todo.app"
	isActive := true
	createdAt := time.Now()

	// Insert hanya kalau belum ada
	_, err = db.Exec(`
		INSERT INTO users (id, name, email, password, role_id, department_id, is_active, created_at, updated_at, deleted_at)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, NULL
		WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = $10)
	`, id, name, email, hashedPassword, roleID, deptID, isActive, createdAt, createdAt, email)

	if err != nil {
		return fmt.Errorf("failed insert user: %v", err)
	}

	fmt.Println("✅ User seeder successfully")
	return nil
}

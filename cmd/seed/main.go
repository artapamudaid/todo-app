package main

import (
	"database/sql"
	"fmt"
	"log"

	"todo-app/internal/config"
	"todo-app/internal/seed"

	_ "github.com/lib/pq"
)

func main() {
	v := config.NewViper()

	username := v.GetString("DATABASE_USERNAME")
	password := v.GetString("DATABASE_PASSWORD")
	host := v.GetString("DATABASE_HOST")
	port := v.GetInt("DATABASE_PORT")
	database := v.GetString("DATABASE_NAME")

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		username, password, database, host, port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("🚀 Running seeders...")

	if err := seed.SeedRoles(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedDepartments(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedUsers(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedProjects(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedBoards(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedCards(db); err != nil {
		log.Fatal(err)
	}
	if err := seed.SeedProjectUsers(db); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ All seeders executed successfully")
}

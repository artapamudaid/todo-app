# Golang Clean Architecture Template

## Description

 Golang Todo App using auth multi user, multi role, multi projects, multi board and PostgreSQL as database

## Architecture

![Clean Architecture](architecture.png)

1. External system perform request (HTTP, gRPC, Messaging, etc)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system 
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

## Tech Stack

- Golang : https://github.com/golang/go
- Posgresql (Database) : https://github.com/postgres/postgres
- Apache Kafka : https://github.com/apache/kafka

## Framework & Library

- GoFiber (HTTP Framework) : https://github.com/gofiber/fiber
- GORM (ORM) : https://github.com/go-gorm/gorm
- Viper (Configuration) : https://github.com/spf13/viper
- Golang Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Logrus (Logger) : https://github.com/sirupsen/logrus
- Sarama (Kafka Client) : https://github.com/IBM/sarama
- JWT (JWT Auth) : github.com/golang-jwt/jwt/v5
- Google UUID : github.com/google/uuid
- Go postgres driver for Go's database/sql package : https://github.com/lib/pq

## Configuration

All configuration is in `.env` file.
Copy `.env-example` to `.env` or you can rename `.env-example` directly to `.env`

## Database Migration

All database migration is in `db/migrations` folder.

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "postgres://username:password@localhost:5432/database?sslmode=disable" -path db/migrations up
```

## Database Seeder

All database seeder is in `internal/seed` folder.

### Seeder Configuration

Seeder configuration is in `cmd/seed/main.go` folder.

### Run Seeder

```bash
go run cmd/seed/main.go
```

## Run Application

### Run web server

```bash
go run cmd/web/main.go
```

### Run worker

```bash
go run cmd/worker/main.go
```
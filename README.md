# Routinist

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile commands)

## Project Structure

```
routinist-go/
├── cmd/
│   └── main.go
├── internal/
│   ├── app/
│   ├── controller/
│   ├── domain/
│   ├── dto/
│   ├── middleware/
│   ├── repository/
│   └── usecase/
├── migrate/
│   └── migrations.sql
├── pkg/
│   └── logger/
├── repo/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
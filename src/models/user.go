package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const USER_IDENTIFIER_STRING string = "user_table"

func CreateUserTable(db *pgx.Conn) {
	log.Printf("Creating [%s] table in database...\n", USER_IDENTIFIER_STRING)

	sqlFile, err := os.ReadFile("src/models/user_create_table.sql")
	if (err != nil) {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	sqlString := fmt.Sprintf(string(sqlFile), USER_IDENTIFIER_STRING)

	sqlResult, err := db.Exec(context.Background(), sqlString)
	if (err != nil) {
		log.Fatalf(
			"Failed to execute statement to create [%s] table in database: %v\n",
			USER_IDENTIFIER_STRING, err,
		)
	}

	log.Printf(
		"Result of create [%s] table in database: %v\n",
		USER_IDENTIFIER_STRING, sqlResult,
	)

	log.Printf("Done creating [%s] table in database.\n", USER_IDENTIFIER_STRING)
}

type User interface {
	// Creates a new user on the provided database, if it doesn't already exist.
	CreateOnDb(db *sql.DB) error
}

type user struct {
	id string
	created_at time.Time
	updated_at time.Time
	name_first string
	name_last string
	email string
}

func (user *user) CreateOnDb(db *sql.DB) error {
	// TODO: Create user on DB
	return nil
}

func NewUser() User {
	// TODO: Initialize User model
	return nil
}

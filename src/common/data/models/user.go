package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

const USER_IDENTIFIER_STRING string = "user_table"

func CreateUserTable(db *pgx.Conn) {
	log.Printf("Creating [%s] table in database...\n", USER_IDENTIFIER_STRING)

	sqlFile, err := os.ReadFile("src/common/data/models/user_create_table.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	sqlString := fmt.Sprintf(string(sqlFile), USER_IDENTIFIER_STRING)

	sqlResult, err := db.Exec(context.Background(), sqlString)
	if err != nil {
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

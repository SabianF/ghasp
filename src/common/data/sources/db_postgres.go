package db_postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitDb() *pgx.Conn {
	log.Println("Opening DB connection...")

	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("Failed to get DB_HOST from .env")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		log.Fatal("Failed to get DB_PORT from .env")
	}

	user := os.Getenv("DB_USER")
	if host == "" {
		log.Fatal("Failed to get DB_USER from .env")
	}

	pass := os.Getenv("DB_PASS")
	if host == "" {
		log.Fatal("Failed to get DB_PASS from .env")
	}

	name := os.Getenv("DB_NAME")
	if host == "" {
		log.Fatal("Failed to get DB_NAME from .env")
	}

	dataSource := fmt.Sprintf(
		// The '' around password is to include any spaces
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
		host, port, user, pass, name,
	)

	db, err := pgx.Connect(context.Background(), dataSource)
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to establish connection to DB: ", err)
	}

	log.Printf("Successfully opened DB connection: %s\n", db.Config().Database)

	return db
}

func CloseDb(db *pgx.Conn) {
	db.Close(context.Background())
}

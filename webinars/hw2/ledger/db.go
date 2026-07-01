package ledger

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func InitDB() error {
	var err error

	dsn := buildDSN()

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = conn
	log.Println("connected to PostgreSQL")
	return nil
}

func buildDSN() string {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		return databaseURL
	}

	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	user := getenvDefault("DB_USER", "postgres")
	password := getenvDefault("DB_PASS", "postgres")
	dbname := getenvDefault("DB_NAME", "cashapp")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func getenvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func CloseDB() error {
	if db == nil {
		return nil
	}

	return db.Close()
}

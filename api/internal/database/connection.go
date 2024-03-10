package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbConn *pgxpool.Pool

func Init() {
	db, err := pgxpool.ConnectConfig(context.Background(), Config())

	if err != nil {
		panic(err)
	}

	dbConn = db
}

func GetConnection() (*pgxpool.Conn, error) {
	return dbConn.Acquire(context.Background())
}

func ReleaseConnection(conn *pgxpool.Conn) {
	conn.Release()
}

func Config() *pgxpool.Config {
	const defaultMaxConns = 50
	const defaultMinConns = 5

	host := os.Getenv("DB_HOST")

	if host == "" {
		host = "localhost"
	}

	// Your own Database URL
	DATABASE_URL := fmt.Sprintf("host=%s port=5432 user=admin password=admin123 dbname=financial sslmode=disable", host)

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.ConnConfig.ConnectTimeout = 30 * time.Second

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before Acquire connection")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("Before Release connection")

		return true
	}

	return dbConfig
}

package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var dbConn *sql.DB

func Init() {

	host := os.Getenv("DB_HOST")

	if host == "" {
		host = "localhost"
	}

	fmt.Println(host)
	config := fmt.Sprintf("host=%s port=5432 user=admin password=admin123 dbname=financial sslmode=disable", host)

	db, err := sql.Open("postgres", config)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	dbConn = db
}

func GetConnection() *sql.DB {
	return dbConn
}

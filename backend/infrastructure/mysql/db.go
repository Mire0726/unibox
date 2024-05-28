package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB(dataSourceName string) *sql.DB {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping to the database: %v", err)
	}

	return db
}

const driverName = "mysql"

var Conn *sql.DB

func ConnectToDB() (*sql.DB, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	charset := os.Getenv("MYSQL_CHARSET")
	parseTime := os.Getenv("MYSQL_PARSE_TIME")
	loc := os.Getenv("MYSQL_LOC")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		user, password, host, port, database, charset, parseTime, loc)

	Conn, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal("cannnot sql.Open", err)
	}
	if err := Conn.Ping(); err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}

	return Conn, nil

}

package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func createConnectionString() string {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func InitDB() {

	connectionString := createConnectionString()
	var err error //required. since := in the next line cause error in creating DB tables.
	DB, err = sql.Open("postgres", connectionString)
	fmt.Println("Connecting to database", connectionString)
	if err != nil {
		panic("Unable to connect to the database1")
	}

	err = DB.Ping()
	if err != nil {
		panic("Unable to connect to the database2")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createTodosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		list JSONB,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		user_id INTEGER
	)
	`
	_, err := DB.Exec(createTodosTable)
	if err != nil {
		panic("Unable to create todos table")
	}
}

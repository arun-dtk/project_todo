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
	fmt.Println(connectionString)
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

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		fmt.Println(err)
		panic("Unable to create users table")
	}

	createTodosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		list JSONB,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createTodosTable)
	if err != nil {
		fmt.Println(err)
		panic("Unable to create todos table")
	}
}

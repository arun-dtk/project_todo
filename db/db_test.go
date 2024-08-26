package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateConnectionString(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")

	connStr := createConnectionString()
	expected := "host=localhost port=5432 user=testuser password=password dbname=testdb sslmode=disable"

	if connStr != expected {
		t.Errorf("Expected connection string %s, but got %s", expected, connStr)
	}
}

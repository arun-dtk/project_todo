package models

import (
	"encoding/json"
	"fmt"
	"project_todo/db"
	"time"
)

type TodoItem struct {
	Item    string `json:"item"`
	Checked bool   `json:"checked"`
}

type Todo struct {
	ID        int64      `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	List      []TodoItem `db:"list" json:"list"`
	IsActive  bool       `db:"is_active" json:"isActive"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	UserID    int64      `db:"user_id" json:"userId"`
}

func (t *Todo) Save() error {
	// Debugging: Print the content of t.List before marshalling
	fmt.Printf("Content of t.List: %+v\n", t.List)

	listJSON, err := json.Marshal(t.List)
	fmt.Println("listJSON", listJSON)
	if err != nil {
		return err
	}
	// Convert listJSON to string
	listJSONString := string(listJSON)
	fmt.Println("listJSONString", listJSONString)

	query := `
	INSERT INTO todos(title, list, is_active, created_at, updated_at, user_id)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing query:", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(t.Title, listJSONString, t.IsActive, t.CreatedAt, t.UpdatedAt, t.UserID).Scan(&t.ID)
	return err
}

func GetAllTodos() ([]Todo, error) {
	query := "SELECT * FROM todos"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []Todo

	for rows.Next() {
		var todo Todo
		var listJson []byte
		err := rows.Scan(&todo.ID, &todo.Title, &listJson, &todo.IsActive, &todo.CreatedAt, &todo.UpdatedAt, &todo.UserID)
		if err != nil {
			return nil, err
		}
		// Unmarshal the JSONB field into the List slice
		err = json.Unmarshal(listJson, &todo.List)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

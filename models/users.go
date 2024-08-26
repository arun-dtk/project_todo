package models

import (
	"errors"
	"fmt"
	"project_todo/db"
	"project_todo/utils"
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `binding:"required" db:"email" json:"email"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Password  string    `binding:"required" db:"password" json:"password"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(email, first_name, last_name, password, is_active, created_at, updated_at)
	VALUES ($1,$2, $3, $4, $5, $6, $7) RETURNING id`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = db.DB.QueryRow(query, u.Email, u.FirstName, u.LastName, hashedPassword, u.IsActive, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	// Scan should has a destination pointer
	fmt.Println(err)
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password from users where email = $1"
	row := db.DB.QueryRow(query, u.Email)
	var existingPassword string
	err := row.Scan(&u.ID, &existingPassword)
	if err != nil {
		return errors.New("Invalid Credentials")
	}

	isValid := utils.ComparePassword(u.Password, existingPassword)
	if !isValid {
		return errors.New("Invalid Credentials")
	}
	return nil
}

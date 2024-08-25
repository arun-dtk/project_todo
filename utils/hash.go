package utils

import "golang.org/x/crypto/bcrypt"

const hashCost = 14

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

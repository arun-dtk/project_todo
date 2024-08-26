package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  "",
		"userId": "",
		"exp":    time.Now().Add(time.Hour * 4).Unix(),
	})
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}

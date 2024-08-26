package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 4).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	fmt.Println("before parsing")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	fmt.Println("after parsing", parsedToken)

	if err != nil {
		return 0, errors.New("Could not parse token")
	}
	isTokenValid := parsedToken.Valid
	fmt.Println("isTokenValid", isTokenValid)

	if !isTokenValid {
		return 0, errors.New("Invalid token!")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	fmt.Println("claims", claims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}
	fmt.Println("UserId.... ", claims["userId"].(float64))
	userId := int64(claims["userId"].(float64))
	return userId, nil
}

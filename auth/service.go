package auth

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/joho/godotenv"
)

var SECRET_KEY []byte
var TOKEN_HOURS int

func init() {
	if err := godotenv.Load("./config/.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
	if len(SECRET_KEY) == 0 {
		log.Fatal("SECRET_KEY is not set in .env file")
	}

	hours, err := strconv.Atoi(os.Getenv("SECRET_TOKEN_HOURS"))
	if err != nil {
		log.Fatal("Invalid SECRET_TOKEN_HOURS value:", err)
	}
	TOKEN_HOURS = hours
}

// GenerateToken generates a JWT token with userID, isAdmin, issued time, and expiry time as payload
func GenerateToken(userID int, isAdmin bool) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userID,
		"isAdmin": isAdmin,
		"iat":     jwt.NewNumericDate(time.Now()),
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Duration(TOKEN_HOURS) * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates the given JWT token and returns the validated token
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		encodedToken,
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid token")
			}
			return SECRET_KEY, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

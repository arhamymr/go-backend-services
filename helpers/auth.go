package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateGlobalToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        "global",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_GLOBAL_TOKEN_KEY")))
	return tokenString
}

func CreateToken(uuid string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        uuid,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_TOKEN_KEY")))
	return tokenString
}

func VerifyToken(tokenString string, tokenEnv string) error {
	// Verify token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(tokenEnv)), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

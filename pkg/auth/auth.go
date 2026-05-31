package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getPasswordHash() string {
	pass := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(hash[:])
}

func GenerateToken() (string, error) {
	secret := os.Getenv("TODO_PASSWORD")
	if secret == "" {
		return "", errors.New("пароль не задан в переменных окружения")
	}

	claims := jwt.MapClaims{
		"hash": getPasswordHash(),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string) bool {
	secret := os.Getenv("TODO_PASSWORD")
	if secret == "" {
		return false
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		tokenHash, okHash := claims["hash"].(string)
		if okHash && tokenHash == getPasswordHash() {
			return true
		}
	}

	return false
}

package jwttoken

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	hmacSampleSecret = "my_secret_key_for_jwt_token"
)

var ErrInvalidToken = errors.New("некорректный токен")

func GenerateToken(user string, duration uint64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user,
		"nbf":  now.Unix(),
		"exp":  now.Add(time.Duration(duration) * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		slog.Error("Невозможно создать токен:", "ОШИБКА:", err)
		return "", ErrInvalidToken
	}
	return tokenString, nil
}

func ParseToken(token string) (string, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		slog.Error("Невозможно распарсить токен:", "ОШИБКА:", err)
		return "", ErrInvalidToken
	}
	if jwtToken.Valid {
		return jwtToken.Claims.(jwt.MapClaims)["name"].(string), nil
	}
	return "", ErrInvalidToken
}

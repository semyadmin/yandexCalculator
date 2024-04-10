package jwttoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	hmacSampleSecret = "my_secret_key_for_jwt_token"
	errInvalidToken  = "invalid token"
)

func GenerateToken(user string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
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
		return "", err
	}
	if jwtToken.Valid {
		return jwtToken.Claims.(jwt.MapClaims)["name"].(string), nil
	}
	return "", errors.New(errInvalidToken)
}

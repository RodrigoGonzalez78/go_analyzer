package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/RodrigoGonzalez78/go_analyzer/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userName string) (string, error) {
	secretKey := []byte("NoMeHackeen")

	claims := models.Claim{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ProcessToken(tokenStr string) (*models.Claim, bool, error) {
	secretKey := []byte("NoMeHackeen")

	if !strings.HasPrefix(tokenStr, "Bearer ") {
		return nil, false, errors.New("formato de token inválido")
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	var claims models.Claim

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, false, err
	}

	if !token.Valid {
		return nil, false, errors.New("token inválido")
	}

	return &claims, true, nil
}

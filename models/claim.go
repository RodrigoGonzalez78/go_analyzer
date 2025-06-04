package models

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

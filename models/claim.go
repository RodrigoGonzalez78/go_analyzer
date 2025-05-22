package models

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

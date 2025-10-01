package security

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID int `json:"sub"`
	jwt.RegisteredClaims
}

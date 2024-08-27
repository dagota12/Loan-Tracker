package domain

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	IsOwner bool   `json:"is_owner"`
	jwt.RegisteredClaims
}
type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

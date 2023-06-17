package models

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	LoginAt  string `json:"loginAt"`
	ExpireAt string `json:"expireAt"`
	jwt.StandardClaims
}

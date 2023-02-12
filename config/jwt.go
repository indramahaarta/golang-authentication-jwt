package config

import "github.com/golang-jwt/jwt/v4"

var JWK_KEY = []byte("aoishdiabdhbachbajkfdcklaj")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
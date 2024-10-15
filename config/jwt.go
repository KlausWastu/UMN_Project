package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("KsjduUg26283hsdSkwQlwK9kS")

type JWTClaim struct {
	UserID int64
	Username string
	jwt.RegisteredClaims
}

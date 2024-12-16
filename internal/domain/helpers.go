package domain

import "github.com/dgrijalva/jwt-go"

type Claim struct {
	Sum  uint   `json:"sum"`
	Role string `json:"role"`
	jwt.StandardClaims
}

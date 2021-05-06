package types

import (
	"github.com/dgrijalva/jwt-go"
)

type Health struct {
	Status string `json:"status"`
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope    string `json:"scope"`
	Email    string `json:"name"`
	UserID   string `json:"sub"`
	NickName string `json:"nickname"`
	jwt.StandardClaims
}

type Url struct {
	Url string `json:"url"`
}

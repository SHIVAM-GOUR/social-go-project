package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuntenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuntenticator(secret, aud, iss string) *JWTAuntenticator {
	return &JWTAuntenticator{secret, iss, aud}
}

func (a *JWTAuntenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuntenticator) ValidateToken(token string) (*jwt.Token, error) {
	return nil, nil
}

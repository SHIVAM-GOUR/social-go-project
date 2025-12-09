package auth

import (
	"fmt"

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
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Name}),
	)
}

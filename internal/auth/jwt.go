package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenIssuer(secret string) *TokenIssuer {
	return &TokenIssuer{secret: []byte(secret), ttl: 7 * 24 * time.Hour}
}

type claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (t *TokenIssuer) Issue(userID int64) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(t.secret)
}

func (t *TokenIssuer) Parse(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}
	c, ok := token.Claims.(*claims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return c.UserID, nil
}

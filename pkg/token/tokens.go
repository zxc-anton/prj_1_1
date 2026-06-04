package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("anton")

type TokenClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func CreateTokens(id int) (string, string, error) {
	accessClaims := TokenClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshClaims := TokenClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accToken, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	refToken, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accToken, refToken, nil
}

func DecodeToken(token string) (int, error) {
	claims := &TokenClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("not valid token")
		}
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	if !jwtToken.Valid {
		return 0, errors.New("not valid token")
	}

	return claims.Id, nil
}

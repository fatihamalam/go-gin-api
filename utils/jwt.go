package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(GetEnv("JWT_SECRET", "your_secret_key"))
var refreshKey = []byte(GetEnv("JWT_SECRET", "your_refresh_key"))

type JWTClaim struct {
	UserID uint `json:"user_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string) (*string, *string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaim{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaim{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	accessToken, err := access.SignedString(jwtKey)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := refresh.SignedString(refreshKey)
	if err != nil {
		return nil, nil, err
	}
	return &accessToken, &refreshToken, nil
}

func ParseToken(tokenString string, isRefresh bool) (*jwt.Token, error) {
	var key []byte
	key = jwtKey
	if isRefresh {
		key = refreshKey
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}

package auth

import (
	"errors"
	"fmt"
	"instancer/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Verify(tokenString string) (*Claims, error) {
	c := env.Get()

	// Define the function to parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and the correct algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and not expired
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {

		d, err := token.Claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}

		if d.Before(time.Now()) {
			return nil, errors.New("expired token")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

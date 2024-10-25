package auth

import (
	"instancer/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	InstanceID string `json:"instanceid"`
	jwt.RegisteredClaims
}

func Generate(instanceID string) (string, error) {
	c := env.Get()

	signingKey := []byte(c.SigningKey)

	// Set custom claims with a 2-day expiry
	claims := Claims{
		InstanceID: instanceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Sign the token with the signing key
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

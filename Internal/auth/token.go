package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadjwt() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("unable to get jwt")
	}
	if len(secret) < 32 {
		return nil, errors.New("token must be more than 32")
	}
	return []byte(secret), nil
}

type Claims struct {
	Role Role `json:"role"`
	jwt.RegisteredClaims
}

func NewToken(secret []byte, accountID string, role Role, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   accountID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)

}

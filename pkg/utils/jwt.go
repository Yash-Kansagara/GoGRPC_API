package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string, userid string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiray := os.Getenv("JWT_EXPIRY")

	duration, err := time.ParseDuration(expiray)
	if err != nil {
		duration = time.Duration(15 * time.Minute)
	}
	claims := jwt.MapClaims{
		"username": username,
		"userid":   userid,
		"role":     role,
		"exp":      time.Now().Add(duration).UTC().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

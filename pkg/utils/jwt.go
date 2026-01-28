package utils

import (
	"log"
	"os"
	"time"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetDefaultRefreshTokenExpiry() time.Duration {
	expiray := os.Getenv("JWT_REFRESH_EXPIRY")
	duration, err := time.ParseDuration(expiray)
	if err != nil {
		log.Println("Failed to parse duration", err)
		duration = time.Duration(24 * 30 * time.Hour)
	}
	return duration
}

func GenerateAccessToken(username string, userid string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiray := os.Getenv("JWT_EXPIRY")

	duration, err := time.ParseDuration(expiray)
	if err != nil {
		log.Println("Failed to parse duration", err)
		duration = time.Duration(15 * time.Minute)
	}
	claims := jwt.MapClaims{
		"username": username,
		"userid":   userid,
		"role":     role,
		"exp":      jwt.NewNumericDate(time.Now().Add(duration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateRefreshToken(username string, userid string, role string) (string, *jwt.Token, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	duration := GetDefaultRefreshTokenExpiry()

	claims := models.RefreshTokenClaims{
		Username: username,
		UserId:   userid,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, token, nil
}

func ParseRefreshToken(token string) (*models.RefreshTokenClaims, *jwt.Token, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	refreshToken := &models.RefreshTokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println("Failed to parse token", err)
		return nil, nil, err
	}
	return refreshToken, parsedToken, nil
}

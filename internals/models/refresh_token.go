package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenClaims struct {
	Username string `json:"username"`
	UserId   string `json:"userid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

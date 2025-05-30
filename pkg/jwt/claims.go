package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/juancanchi/users/internal/domain"
	"time"
)

func GenerateJWT(user *domain.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
		"role":    user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

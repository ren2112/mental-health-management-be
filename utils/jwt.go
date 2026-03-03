package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var JwtSecret = []byte("your-secret-key") // 建议放配置文件

type CustomClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"` // ⭐ 直接 int

	jwt.RegisteredClaims
}

func GenerateJWT(userID int, role int) (string, error) {

	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(7 * 24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "mental-health-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

func ParseJWT(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

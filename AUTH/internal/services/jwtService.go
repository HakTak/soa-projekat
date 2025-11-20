package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secret, issuer string) *JWTService {
	return &JWTService{
		secretKey: secret,
		issuer:    issuer,
	}
}

func (j *JWTService) GenerateToken(userId string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"iss":     j.issuer,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(j.secretKey), nil
	})
}

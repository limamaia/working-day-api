package helpers

import (
	"fmt"
	"time"
	"working-day-api/config"
	"working-day-api/internal/domain"

	"github.com/dgrijalva/jwt-go"
)

func NewJWTService() *jwtService {
	return &jwtService{
		SecretKey: config.AppConfig.SecretKey,
		Issuer:    config.AppConfig.Issuer,
	}
}

type jwtService struct {
	SecretKey string
	Issuer    string
}

func (s *jwtService) GenerateToken(id uint, role string) (string, error) {
	claim := &domain.Claim{
		Sum:  id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    s.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) ParseToken(token string) (*domain.Claim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &domain.Claim{}, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*domain.Claim); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

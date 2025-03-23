package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"

	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/models"
)

func GenerateJWTToken(user *models.User, ttl int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": user.Role,
		"exp":  time.Now().Add(time.Duration(ttl) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(length int) (string, error) {
	if length <= 0 {
		return "", appError.ErrInvalidTokenLength
	}

	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", appError.ErrFailedRandGeneration
	}

	tokenString := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(tokenBytes)
	return tokenString, nil
}

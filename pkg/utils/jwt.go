package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

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

func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appError.ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, appError.ErrInvalidToken
	}

	return claims, nil
}

func extractUserID(claims jwt.MapClaims) (uuid.UUID, error) {
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, appError.ErrInvalidToken
	}

	return uuid.Parse(userIDStr)
}

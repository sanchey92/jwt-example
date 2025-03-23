package utils

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/models"
)

const (
	testSecret = "testSecret"
	testTTL    = 5 // minutes
)

func TestGenerateJWTToken(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		secret  string
		ttl     int
		wantErr bool
	}{
		{
			name: "valid token generation",
			user: &models.User{
				ID:   uuid.New(),
				Role: models.RoleUser,
			},
			secret:  testSecret,
			ttl:     testTTL,
			wantErr: false,
		},
		{
			name: "expired token",
			user: &models.User{
				ID:   uuid.New(),
				Role: models.RoleUser,
			},
			ttl:     -1, // Token already expired
			secret:  testSecret,
			wantErr: false, // Generation succeeds, but token will be invalid on verification
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := GenerateJWTToken(tt.user, tt.ttl, tt.secret)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, tokenStr)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, tokenStr)

			// Validate token
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.secret), nil
			})

			if tt.name == "expired token" {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, token.Valid)

			claims, ok := token.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, tt.user.ID.String(), claims["sub"])
			assert.Equal(t, string(tt.user.Role), claims["role"])

			exp, ok := claims["exp"].(float64)
			assert.True(t, ok)
			assert.WithinDuration(
				t,
				time.Now().Add(time.Duration(tt.ttl)*time.Minute),
				time.Unix(int64(exp), 0),
				time.Minute,
			)
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		wantErr     bool
		expectedErr error
		checkLength bool
	}{
		{
			name:        "Valid length",
			length:      32,
			wantErr:     false,
			checkLength: true,
		},
		{
			name:        "Zero length",
			length:      0,
			wantErr:     true,
			expectedErr: appError.ErrInvalidTokenLength,
		},
		{
			name:        "Negative length",
			length:      -1,
			wantErr:     true,
			expectedErr: appError.ErrInvalidTokenLength,
		},
		{
			name:        "Short length",
			length:      16,
			wantErr:     false,
			checkLength: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tt.length)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
				assert.NoError(t, err)
				assert.NotEmpty(t, decoded)

				if tt.checkLength {
					// Calculate expected length after base64 encoding
					// For URL-safe base64 without padding: ceil(n * 4/3)
					expectedLength := ((tt.length * 4) + 2) / 3
					assert.Equal(t, expectedLength, len(token))
				}

				token2, err := GenerateRefreshToken(tt.length)
				assert.NoError(t, err)
				assert.NotEqual(t, token, token2)
			}
		})
	}
}

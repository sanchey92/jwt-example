package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sanchey92/jwt-example/internal/config"
	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/service"
	"github.com/sanchey92/jwt-example/pkg/utils"
)

func Authenticate(service *service.AuthService, cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, appError.Unauthorized(appError.ErrUnauthorized))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			user, err := service.ExtractUserFromToken(r.Context(), tokenStr, cfg.JWTAccessSecret)
			if err != nil {
				if errors.Is(err, appError.ErrTokenExpired) {
					handleTokenExpired(w, r, service, cfg, next)
					return
				}
				writeError(w, appError.Unauthorized(err))
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleTokenExpired(w http.ResponseWriter, r *http.Request, service *service.AuthService, cfg *config.Config, next http.Handler) {
	tokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeError(w, appError.Unauthorized(appError.ErrUnauthorized))
		return
	}

	refreshToken, user, err := service.ExtractUserFromRefreshToken(r.Context(), tokenCookie.Value)
	if err != nil {
		writeError(w, appError.Unauthorized(err))
		return
	}

	if service.IsRefreshTokenExpired(refreshToken) {
		newRefresh, err := service.GetNewRefreshToken(r.Context(), user.ID, refreshToken.Token)
		if err != nil {
			writeError(w, appError.Unauthorized(err))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefresh,
			HttpOnly: true,
		})
	}

	newAccess, err := utils.GenerateJWTToken(user, cfg.AccessTokenTTL, cfg.JWTAccessSecret)
	if err != nil {
		writeError(w, appError.Unauthorized(appError.ErrInternalServer))
		return
	}

	w.Header().Set("Authorization", "Bearer "+newAccess)

	ctx := context.WithValue(r.Context(), "user", user)

	next.ServeHTTP(w, r.WithContext(ctx))
}

func writeError(w http.ResponseWriter, apiErr *appError.ApiError) {
	w.WriteHeader(apiErr.StatusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": apiErr.Message})
}

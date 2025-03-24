package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/logger"
	"github.com/sanchey92/jwt-example/internal/models"
)

const (
	MaxRequestSize = 1048576 // 1MB
)

type AuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	service   AuthService
	log       *zap.Logger
	validator *validator.Validate
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service:   service,
		log:       logger.GetLogger(),
		validator: validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input AuthInput

	if err := h.decodeJSON(w, r, &input); err != nil {
		h.log.Error("Decoding JSON error", zap.Error(err))
		return
	}

	user, err := h.service.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		h.log.Error("Registration error", zap.Error(err), zap.String("email", input.Email))
		h.writeError(w, appError.InternalServer(err))
		return
	}

	h.log.Info("Successful registration", zap.String("email", user.Email))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input AuthInput

	if err := h.decodeJSON(w, r, &input); err != nil {
		h.log.Error("Decoding JSON error", zap.Error(err))
		return
	}

	tokenPair, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		h.log.Error("Login error", zap.Error(err), zap.String("email", input.Email))
		h.writeError(w, appError.Unauthorized(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		HttpOnly: true,
	})

	h.log.Info("success login", zap.String("email", input.Email))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenPair.AccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input RefreshTokenInput

	if err := h.decodeJSON(w, r, &input); err != nil {
		h.log.Error("Decoding JSON error", zap.Error(err))
		return
	}

	if err := h.service.Logout(r.Context(), input.RefreshToken); err != nil {
		h.log.Error("Logout error", zap.Error(err))
		h.writeError(w, appError.InternalServer(err))
		return
	}

	// Add access token to the black list (with redis for example)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		h.writeError(w, appError.Unauthorized(appError.ErrUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		h.writeError(w, appError.BadRequest(appError.ErrInvalidInput))
		return appError.ErrInvalidInput
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.writeError(w, appError.BadRequest(appError.ErrInvalidInput))
		return appError.ErrInvalidInput
	}

	if err := h.validator.Struct(v); err != nil {
		h.writeError(w, appError.BadRequest(err))
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (h *AuthHandler) writeError(w http.ResponseWriter, apiError *appError.ApiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiError.StatusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": apiError.Message})
}

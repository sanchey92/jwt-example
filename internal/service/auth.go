package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/sanchey92/jwt-example/internal/config"
	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/logger"
	"github.com/sanchey92/jwt-example/internal/models"
	"github.com/sanchey92/jwt-example/pkg/utils"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type TokenRepository interface {
	SaveToken(ctx context.Context, token *models.RefreshToken) error
	GetToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteToken(ctx context.Context, token string) error
}

type AuthService struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
	cfg       *config.Config
	log       *zap.Logger
}

func NewAuthService(userRepo UserRepository, tokenRepo TokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
		log:       logger.GetLogger(),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to get hashed password", zap.Error(err))
		return nil, appError.InternalServer(err)
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		Role:      models.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("Failed to save new user to database", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, appError.ErrUserNotFound) {
			return nil, appError.Unauthorized(appError.ErrUserNotFound)
		}
		return nil, appError.InternalServer(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, appError.Unauthorized(appError.ErrInvalidPassword)
	}

	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, appError.InternalServer(err)
	}

	if err = s.saveRefreshToken(ctx, user.ID, tokenPair.RefreshToken); err != nil {
		return nil, appError.InternalServer(err)
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.DeleteToken(ctx, refreshToken)
}

func (s *AuthService) ExtractUserFromToken(ctx context.Context, tokenStr, secret string) (*models.User, error) {
	claims, err := utils.ParseToken(tokenStr, secret)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ExtractUserID(claims)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, appError.Unauthorized(err)
	}

	return user, nil
}

func (s *AuthService) generateTokenPair(user *models.User) (*models.TokenPair, error) {
	accessToken, err := utils.GenerateJWTToken(user, s.cfg.AccessTokenTTL, s.cfg.JWTAccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) saveRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	refreshToken := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.RefreshTokenTTL) * 24 * time.Hour),
	}

	if err := s.tokenRepo.SaveToken(ctx, refreshToken); err != nil {
		return err
	}
	return nil
}

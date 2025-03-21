package pg

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/logger"
	"github.com/sanchey92/jwt-example/internal/models"
)

type Storage struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	log := logger.GetLogger()

	if err := ctx.Err(); err != nil {
		log.Error("context canceled before connecting to postgres", zap.Error(err))
		return nil, err
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("failed to parse postgres config", zap.Error(err))
		return nil, err
	}

	// TODO: add config (env) for database
	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("failed to create new postgres pool", zap.Error(err))
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Error("failed to ping postgres db", zap.Error(err))
		pool.Close()
		return nil, err
	}

	log.Info("success connect to database")
	return &Storage{
		db:  pool,
		log: log,
	}, nil
}

func (s *Storage) DB() *pgxpool.Pool {
	if s.db != nil {
		return s.db
	}

	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		s.db.Close()
		s.log.Info("Close connection to database")
	}
	return nil
}

func (s *Storage) Create(ctx context.Context, user *models.User) error {
	_, err := s.db.Exec(ctx, createUser, user.ID, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return appError.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *Storage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx, findByEmail, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx, findById, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) SaveToken(ctx context.Context, token *models.RefreshToken) error {
	_, err := s.db.Exec(ctx, saveToken, token.ID, token.UserID, token.Token, token.ExpiresAt)
	return err
}

func (s *Storage) GetToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := s.db.QueryRow(ctx, getToken, token).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrInvalidToken
		}
		return nil, err
	}

	return &t, nil
}

func (s *Storage) DeleteToken(ctx context.Context, token string) error {
	_, err := s.db.Exec(ctx, deleteToken, token)
	return err
}

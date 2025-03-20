package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/sanchey92/jwt-example/internal/logger"
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

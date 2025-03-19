package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	PgDSN            string
	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTokenTTL   int // minute
	RefreshTokenTTL  int // days
}

func MustLoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:             os.Getenv("PORT"),
		PgDSN:            os.Getenv("PG_DSN"),
		JWTAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
	}

	if cfg.Port == "" || cfg.PgDSN == "" || cfg.JWTAccessSecret == "" || cfg.JWTRefreshSecret == "" {
		panic("Failed to get env variables")
	}

	cfg.AccessTokenTTL = 15 // 15 minutes
	cfg.RefreshTokenTTL = 7 // 7 days

	return cfg
}

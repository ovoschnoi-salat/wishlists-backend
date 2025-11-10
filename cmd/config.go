package main

import (
	"backend/internal/config"
	"backend/pkg/http"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type DatabaseConfig struct {
	DSN string `env:"PG_DSN,required"`
}

type Bot struct {
	Token     string `env:"BOT_TOKEN"`
	WebAppURL string `env:"TELEGRAM_WEBAPP_URL"`
}

type Config struct {
	Stage      config.Stage `env:"STAGE" envDefault:"dev"`
	Database   DatabaseConfig
	Bot        Bot
	HttpServer http.ServerConfig
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Err(err).Msg("error loading .env file")
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return &cfg, nil
}

package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"link-service/internal/logger"
	"link-service/internal/server"
)

type Config struct {
	HTTPServer server.Config
	Logger     logger.Config
}

func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

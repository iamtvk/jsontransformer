package config

import (
	"os"
	"time"
)

type Config struct {
	DefaultTransformTimeout time.Duration
	DbUrl                   string
	HTTPPort                string
	GRPCPort                string
}

func Load() *Config {
	return &Config{
		DefaultTransformTimeout: 2 * time.Second,
		DbUrl:                   os.Getenv("DB_URL"),
		HTTPPort:                os.Getenv("HTTP_PORT"),
		GRPCPort:                os.Getenv("GRPC_PORT"),
	}
}

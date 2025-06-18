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
		DbUrl:                   "postgres://postgres:password@localhost:5432/myapp?sslmode=disable",
		HTTPPort:                os.Getenv("HTTP_PORT"),
		GRPCPort:                os.Getenv("GRPC_PORT"),
	}
}

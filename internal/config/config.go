package config

import "time"

type Config struct {
	DefaultTimeout time.Duration
	DbUrl          string
}

func Load() *Config {
	return &Config{
		DefaultTimeout: time.Second,
		DbUrl:          "localhost:6969",
	}
}

package config

import "os"

type Config struct {
	RedisAddr string
}

func Load() *Config {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // default for local dev
	}

	return &Config{
		RedisAddr: addr,
	}
}

package config

import (
	"time"
)

type Config struct {
	Addr           string
	HeatbeatTTL    time.Duration
	WorkerPoolSize int
}

func LoadConfig() *Config {
	return &Config{
		Addr:           ":12345",
		HeatbeatTTL:    60 * time.Second,
		WorkerPoolSize: 10,
	}
}

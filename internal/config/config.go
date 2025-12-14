package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	GRPCAddr string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	cfg := &Config{
		HTTPAddr: getEnv("HTTPPort", "8080"),
		GRPCAddr: getEnv("GRPCPort", "9091"),
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return i
}

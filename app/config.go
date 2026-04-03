package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV           string
	PORT          string
	DATABASE_HOST string
	DATABASE_PORT string
	DATABASE_NAME string
	DATABASE_USER string
	DATABASE_PASS string
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.DATABASE_USER, c.DATABASE_PASS, c.DATABASE_HOST, c.DATABASE_PORT, c.DATABASE_NAME,
	)
}

func loadConfig() *Config {
	godotenv.Load()

	return &Config{
		ENV:           getEnv("ENV", "development"),
		PORT:          getEnv("PORT", "8080"),
		DATABASE_HOST: getEnv("DATABASE_HOST", "localhost"),
		DATABASE_PORT: getEnv("DATABASE_PORT", "5432"),
		DATABASE_NAME: mustGetEnv("DATABASE_NAME"),
		DATABASE_USER: mustGetEnv("DATABASE_USER"),
		DATABASE_PASS: mustGetEnv("DATABASE_PASS"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %q is not set", key))
	}
	return v
}

package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// NewLogger returns a JSON logger in production and a colorized pretty logger
// in every other environment (development, test, etc.).
func NewLogger(cfg *Config) *slog.Logger {
	if cfg.ENV == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(tint.NewHandler(os.Stdout, nil))
}

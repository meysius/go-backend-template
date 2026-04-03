package app

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func newLogger(env string) *slog.Logger {
	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(tint.NewHandler(os.Stdout, nil))
}

package logger

import (
	"log/slog"
	"os"
	"sso/internal/config"
	slogpretty "sso/internal/lib/logger/handlers/slogprettty"
)

func SetUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case config.EnvLocal:
		log = setUpPrettySlog()
	case config.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}
	return log
}

func setUpPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

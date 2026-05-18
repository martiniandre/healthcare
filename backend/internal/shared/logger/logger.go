package logger

import (
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
)

func Init(env, sentryDSN string) {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      env,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			slog.Error("Failed to initialize Sentry", "error", err)
		} else {
			slog.Info("Sentry initialized successfully")
		}
	}
}

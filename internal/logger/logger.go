package logger

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})

	l := slog.New(handler)
	slog.SetDefault(l)

	return l
}

func parseLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.LevelInfo
	}

	return l
}

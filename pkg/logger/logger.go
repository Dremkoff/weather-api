package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	writer := os.Stdout
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(writer, options)

	return slog.New(handler)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

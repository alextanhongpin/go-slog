package main

import (
	"os"
	"time"

	"golang.org/x/exp/slog"
)

func main() {
	attrs := []slog.Attr{
		slog.Time("build-date", time.Now()),
		slog.String("app-version", "v1.0.0"),
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout).WithAttrs(attrs)
	_ = jsonHandler

	textHandler := slog.NewTextHandler(os.Stdout).WithAttrs(attrs)
	_ = textHandler

	logger := slog.New(jsonHandler)
	logger.Info("greetings",
		slog.String("msg", "hello"),
		slog.Group("user",
			slog.String("name", "John"),
			slog.Int("age", 10),
		),
	)
	slog.SetDefault(logger)
	slog.Info("hello")
}

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"log/slog"
)

func main() {
	attrs := []slog.Attr{
		slog.Time("build-date", time.Now()),
		slog.String("app-version", "v1.0.0"),
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}).WithAttrs(attrs)
	_ = jsonHandler

	textHandler := slog.NewTextHandler(os.Stdout, nil).WithAttrs(attrs)
	_ = textHandler

	logger := slog.New(&ContextHandler{Handler: jsonHandler})
	logger.Info("greetings",
		slog.String("msg", "hello"),
		slog.Group("user",
			slog.String("name", "John"),
			slog.Int("age", 10),
			Err(errors.New("bad request")),
		),
	)

	ctx := context.Background()
	ctx = Append(ctx, "added", true)
	ctx = Append(ctx, "meaning_of_life", 42)
	ctx = context.WithValue(ctx, "user_id", "123")
	logger.ErrorContext(ctx, "failed to greet", Err(&UserError{ID: "123"}))
	slog.SetDefault(logger)
	slog.Info("hello")
}

type UserError struct {
	ID string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("User(id: %s)", e.ID)
}

func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
	//return slog.Any("err", err)
}

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Add("foo", "bar")
	r.Add("married", true)

	var attrs []slog.Attr
	if userID, ok := ctx.Value("user_id").(string); ok {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if as, ok := ctx.Value("slog_attrs").(*[]slog.Attr); ok {
		attrs = append(attrs, *as...)
	}
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}

func Append(ctx context.Context, key string, value any) context.Context {
	if attrs, ok := ctx.Value("slog_attrs").(*[]slog.Attr); ok {
		*attrs = append(*attrs, slog.Any(key, value))
		return ctx
	}

	return context.WithValue(ctx, "slog_attrs", &[]slog.Attr{slog.Any(key, value)})
}

package main

import (
	"context"
	"log/slog"
)

type key uint8

const (
	loggerKey key = iota
	cfgKey
)

func logger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return l
}

func withLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func cfg(ctx context.Context) *Config {
	return ctx.Value(cfgKey).(*Config)
}

func withConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, cfgKey, cfg)
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"imdb.tsv.analyzer/internal/analyzer"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	slog.InfoContext(ctx, fmt.Sprintf("app started at %s", start.Format(time.RFC3339)))

	cfg, err := analyzer.NewConfig()
	if err != nil {
		slog.ErrorContext(ctx, "error config loading", "err", err)
		return
	}

	app, err := analyzer.NewApp(*cfg)
	if err != nil {
		slog.ErrorContext(ctx, "error app init", "err", err)
		return
	}

	err = app.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error app run", "err", err)
		return
	}

	slog.InfoContext(ctx, "app finished in ", "time", time.Since(start).String())
}

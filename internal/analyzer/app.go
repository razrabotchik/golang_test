package analyzer

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"imdb.tsv.analyzer/internal/pkg/imdbws"
)

type App struct {
	cfg Config
}

func NewApp(cfg Config) (*App, error) {
	return &App{
		cfg: cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(a.cfg.MaxRunTime))
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	opts, err := InitOptions(ctx, a.cfg)
	if err != nil {
		slog.ErrorContext(ctx, "error options init", "err", err)
		return err
	}

	imdb, err := imdbws.NewClient(a.cfg.Imsdbws)
	if err != nil {
		slog.ErrorContext(ctx, "error imdbws init", "err", err)
		return err
	}

	parser, err := NewParser(imdb, opts.Filters, opts.ThreadsCount)
	if err != nil {
		slog.ErrorContext(ctx, "error parser init", "err", err)
		return err
	}

	tsv := NewTsvReader(opts.FilePath)
	var wg sync.WaitGroup

	films := NewFilms()
	collector := make(chan imdbws.Film)

	threads := parser.InitThreads(ctxTimeout, &wg, collector)
	go films.Store(ctxTimeout, collector)

	iteration := 0
	err = tsv.Process(ctxTimeout, func(lineNumber uint, line string) {
		if lineNumber > 1 {
			iteration += 1
			if iteration > (len(threads) - 1) {
				iteration = 0
			}

			threads[iteration] <- line
		}
	})
	if err != nil {
		slog.ErrorContext(ctx, "error file read", "err", err)
		return err
	}

	slog.InfoContext(ctx, "lines count", "count", parser.Count())
	parser.Close(threads)
	close(collector)

	wg.Wait()

	for f := range films.GetFilms() {
		slog.InfoContext(ctx, "film", "film", f)
	}

	slog.InfoContext(ctx, "app finished")
	return nil
}

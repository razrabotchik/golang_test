package analyzer

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
)

type Options struct {
	FilePath string
	Filters  Filters

	MaxRunTime   uint // maximum run time of the application. Format is a `time.Duration` string see [here](https://godoc.org/time#ParseDuration)
	ThreadsCount uint
}

func InitOptions(ctx context.Context, cfg Config) (Options, error) {
	filters := Filters{}

	flag.StringVar(&filters.ID, "id", "", "filter id")
	flag.StringVar(&filters.TitleType, "titleType", "", "filter titleType")
	flag.StringVar(&filters.PrimaryTitle, "primaryTitle", "", "filter primaryTitle")
	flag.StringVar(&filters.OriginalTitle, "originalTitle", "", "filter originalTitle")
	flag.StringVar(&filters.StartYear, "startYear", "", "filter startYear")
	flag.StringVar(&filters.EndYear, "endYear", "", "filter endYear")
	flag.StringVar(&filters.RuntimeMinutes, "runtimeMinutes", "", "filter runtimeMinutes")
	flag.StringVar(&filters.Genre, "genre", "", "filter genre")
	flag.StringVar(&filters.Genres, "genres", "", "filter genres")
	flag.StringVar(&filters.PlotFilter, "plotFilter", "", "filter plotFilter")

	var filePath string
	flag.StringVar(&filePath, "file", "", "path to tsv file")
	flag.Parse()

	if filePath == "" {
		slog.ErrorContext(ctx, "error: --file is empty")
		return Options{}, fmt.Errorf("file is empty")
	}

	options := Options{
		FilePath:     filePath,
		Filters:      filters,
		ThreadsCount: cfg.ThreadsCount,
	}

	flag.UintVar(&options.MaxRunTime, "maxRunTime", 0, "filter maxRunTime")
	if options.MaxRunTime == 0 {
		options.MaxRunTime = cfg.MaxRunTime
	}

	return options, nil
}

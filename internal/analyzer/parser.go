package analyzer

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"imdb.tsv.analyzer/pkg/imdbws"
)

type Imdbws interface {
	GetFilmByID(ctx context.Context, id string) (*imdbws.Film, error)
}

type Filters struct {
	ID             string
	TitleType      string // filter on `titleType` column
	PrimaryTitle   string // filter on `primaryTitle` column
	OriginalTitle  string // filter on `originalTitle` column
	Genre          string // filter on `genre` column
	StartYear      string // filter on `startYear` column
	EndYear        string // filter on `endYear` column
	RuntimeMinutes string // filter on `runtimeMinutes` column
	Genres         string // filter on `genres` column
	PlotFilter     string // regex pattern to apply to the plot of a film retrieved from [omdbapi](https://www.omdbapi.com/)
}

type Parser struct {
	imdb         Imdbws
	filters      Filters
	threadsCount uint
	count        atomic.Uint64
}

func NewParser(imdb Imdbws, filters Filters, threadsCount uint) (*Parser, error) {
	if threadsCount == 0 {
		return nil, errors.New("threadsCount cannot be zero")
	}

	return &Parser{
		imdb:         imdb,
		filters:      filters,
		threadsCount: threadsCount,
		count:        atomic.Uint64{},
	}, nil
}

func filterCheck(filter string, cell string, match int, filtersNeeded int) (int, int) {
	if filter != "" {
		filtersNeeded += 1
		if cell == filter {
			match += 1
		}
	}

	return match, filtersNeeded
}

func (p *Parser) Parse(ctx context.Context, line string) (bool, *imdbws.Film) {
	cells := strings.Split(line, "\t")

	if len(cells) < 9 {
		return false, nil
	}

	filtersNeeded := 0
	match := 0

	match, filtersNeeded = filterCheck(p.filters.ID, cells[0], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.TitleType, cells[1], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.PrimaryTitle, cells[2], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.OriginalTitle, cells[3], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.Genre, cells[4], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.StartYear, cells[5], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.EndYear, cells[6], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.RuntimeMinutes, cells[7], match, filtersNeeded)
	match, filtersNeeded = filterCheck(p.filters.Genres, cells[8], match, filtersNeeded)

	var film *imdbws.Film
	var err error
	if p.filters.PlotFilter != "" && match == filtersNeeded { //to not make request if other filters not equal
		filtersNeeded += 1
		film, err = p.imdb.GetFilmByID(ctx, cells[0])
		if err != nil {
			slog.ErrorContext(ctx, "error get film by id", "error", err)
			return false, nil
		}

		if film != nil {
			m, err := regexp.MatchString(p.filters.PlotFilter, film.Plot)
			if err != nil {
				slog.ErrorContext(ctx, "error regexp match", "error", err)
				return false, nil
			}

			if m {
				match += 1
			}
		}
	}

	if match == filtersNeeded {
		return true, film
	}

	return false, nil
}

func (p *Parser) Analyze(ctx context.Context, wg *sync.WaitGroup, data <-chan string, collector chan<- imdbws.Film) {
	defer wg.Done()

	for {
		select {
		case l, ok := <-data:
			if !ok {
				return
			}

			if ok, film := p.Parse(ctx, l); ok {
				collector <- *film
				slog.DebugContext(ctx, "film found", "line", l)
			}

			p.count.Add(1)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Parser) InitThreads(ctx context.Context, wg *sync.WaitGroup, collector chan imdbws.Film) []chan string {
	// init workers
	threads := make([]chan string, p.threadsCount)
	for k := range threads {
		wg.Add(1)
		thread := make(chan string)
		go p.Analyze(ctx, wg, thread, collector)
		threads[k] = thread
	}

	slog.InfoContext(ctx, "workers count", "count", p.threadsCount)
	return threads
}

func (p *Parser) Count() uint64 {
	return p.count.Load()
}

func (p *Parser) Close(threads []chan string) {
	for k := range threads {
		close(threads[k])
	}
}

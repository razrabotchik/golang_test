package infra

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type Filters struct {
	Id             string
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

type Analyzer struct {
	wg        *sync.WaitGroup
	imdb      *Imdbws
	filters   Filters
	chanCount int
	Count     int
}

func NewAnalyzer(wg *sync.WaitGroup, imdb *Imdbws, filters Filters, chanCount int) Analyzer {
	return Analyzer{
		wg:        wg,
		imdb:      imdb,
		filters:   filters,
		chanCount: chanCount,
		Count:     0,
	}
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

func (a *Analyzer) Parse(ctx context.Context, line string) bool {
	cells := strings.Split(line, "\t")

	if len(cells) < 9 {
		return false
	}

	filtersNeeded := 0
	match := 0

	match, filtersNeeded = filterCheck(a.filters.Id, cells[0], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.TitleType, cells[1], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.PrimaryTitle, cells[2], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.OriginalTitle, cells[3], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.Genre, cells[4], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.StartYear, cells[5], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.EndYear, cells[6], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.RuntimeMinutes, cells[7], match, filtersNeeded)
	match, filtersNeeded = filterCheck(a.filters.Genres, cells[8], match, filtersNeeded)

	if a.filters.PlotFilter != "" && match == filtersNeeded { //to not make request if other filters not equal
		filtersNeeded += 1
		film := a.imdb.GetById(ctx, cells[0])
		if film != nil {
			m, _ := regexp.MatchString(a.filters.PlotFilter, film.Plot)
			if m {
				match += 1
			}
		}
	}

	if match == filtersNeeded {
		return true
	}

	return false
}

func (a *Analyzer) InitChan(ctx context.Context) []chan string {
	// init chan
	lineCans := make([]chan string, a.chanCount)
	for k, _ := range lineCans {
		lineCans[k] = make(chan string)
		go a.Analyze(ctx, lineCans[k], a.wg)
	}

	return lineCans
}

func (a *Analyzer) Analyze(ctx context.Context, data <-chan string, wg *sync.WaitGroup) {
	for {
		select {
		case l := <-data:
			if a.Parse(ctx, l) {
				fmt.Println(l)
			}
			a.Count += 1
			wg.Done()
		case <-ctx.Done():
			return
		}
	}
}

func (a *Analyzer) Close(lineCans []chan string) {
	for k, _ := range lineCans {
		close(lineCans[k])
	}
}

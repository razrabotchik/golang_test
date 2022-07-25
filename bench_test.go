package main

import (
	"context"
	"infoftex/infra"
	"sync"
	"testing"
	"time"
)

func getAnalyzer(filters *infra.Filters) infra.Analyzer {
	var wg sync.WaitGroup
	imdb := infra.NewImdbws("a", 1)
	if filters == nil {
		filters = &infra.Filters{
			Id:             "",
			TitleType:      "",
			PrimaryTitle:   "",
			OriginalTitle:  "",
			Genre:          "",
			StartYear:      "",
			EndYear:        "",
			RuntimeMinutes: "",
			Genres:         "",
			PlotFilter:     "",
		}
	}

	return infra.NewAnalyzer(&wg, &imdb, *filters, 1)
}

func BenchmarkAnalyzeParse(b *testing.B) {
	a := getAnalyzer(nil)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	for i := 0; i < b.N; i++ {
		if r := a.Parse(ctxTimeout, "tt0000014\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short"); r != true {
			b.Fatal("Unexpected response: false")
		}
	}
}

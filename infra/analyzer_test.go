package infra_test

import (
	"context"
	"infoftex/infra"
	"sync"
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	a := getAnalyzer(nil)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := a.Parse(ctxTimeout, "tt0000002\tshort\tLe clown et ses chiens\tLe clown et ses chiens\t0\t1892\t\\N\t5\tAnimation,Short")

	if res != true {
		t.Error("analizer.Analyze is incorrect")
	}
}

func TestAnalyzeIncorrectLine(t *testing.T) {
	a := getAnalyzer(nil)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := a.Parse(ctxTimeout, "short\tLe clown et ses chiens\tLe clown et ses chiens\t0\t1892\t\\N\t5\tAnimation,Short")

	if res != false {
		t.Error("analizer.Analyze incorrect line length checker")
	}
}

func TestAnalyzeFiltering(t *testing.T) {
	filters := infra.Filters{
		Id:             "tt0000014",
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

	a := getAnalyzer(&filters)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := a.Parse(ctxTimeout, "tt0000014\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short")

	if res != true {
		t.Error("analizer.Analyze filter by id is incorrect")
	}

	res = a.Parse(ctxTimeout, "tt0000013\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short")

	if res != false {
		t.Error("analizer.Analyze filter by id is incorrect")
	}
}

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

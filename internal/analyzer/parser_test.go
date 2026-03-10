package analyzer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"imdb.tsv.analyzer/internal/analyzer"
	"imdb.tsv.analyzer/internal/analyzer/mocks"
	"imdb.tsv.analyzer/internal/pkg/imdbws"
)

func Test_Analyze(t *testing.T) {
	a := getAnalyzer(t, nil, mocks.NewImdbws(t))

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, film := a.Parse(ctxTimeout, "tt0000002\tshort\tLe clown et ses chiens\tLe clown et ses chiens\t0\t1892\t\\N\t5\tAnimation,Short")
	assert.True(t, res, "analizer.Analyze is incorrect")
	assert.Nil(t, film, "film should be nil")
}
func Test_AnalyzeWithRequest(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	imdb := mocks.NewImdbws(t)
	imdb.On("GetFilmByID", mock.Anything, "tt0000002").
		Return(&imdbws.Film{Title: "Le clown et ses chiens", Plot: "full"}, nil)

	filters := &analyzer.Filters{
		ID:         "tt0000002",
		PlotFilter: "full",
	}

	a := getAnalyzer(t, filters, imdb)

	res, film := a.Parse(ctxTimeout, "tt0000002\tshort\tLe clown et ses chiens\tLe clown et ses chiens\t0\t1892\t\\N\t5\tAnimation,Short")

	assert.True(t, res, "analizer.Analyze is incorrect")
	assert.NotNil(t, film, "film should not be nil")

	imdb.AssertExpectations(t)
}

func Test_AnalyzeIncorrectLine(t *testing.T) {
	a := getAnalyzer(t, nil, mocks.NewImdbws(t))

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, film := a.Parse(ctxTimeout, "short\tLe clown et ses chiens\tLe clown et ses chiens\t0\t1892\t\\N\t5\tAnimation,Short")
	assert.False(t, res, "analizer.Analyze is incorrect")
	assert.Nil(t, film, "film should be nil")
}

func Test_AnalyzeFiltering(t *testing.T) {
	filters := analyzer.Filters{
		ID:             "tt0000014",
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

	a := getAnalyzer(t, &filters, mocks.NewImdbws(t))

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, film := a.Parse(ctxTimeout, "tt0000014\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short")
	assert.True(t, res, "analizer.Analyze is incorrect")
	assert.Nil(t, film, "film should be nil")

	res, film = a.Parse(ctxTimeout, "tt0000013\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short")
	assert.False(t, res, "analizer.Analyze is incorrect")
	assert.Nil(t, film, "film should be nil")
}

func getAnalyzer(t assert.TestingT, filters *analyzer.Filters, imdb analyzer.Imdbws) *analyzer.Parser {
	if filters == nil {
		filters = &analyzer.Filters{}
	}

	parser, err := analyzer.NewParser(imdb, *filters, 1)
	assert.NoError(t, err)

	return parser
}

func BenchmarkAnalyzeParse(b *testing.B) {
	a := getAnalyzer(b, nil, mocks.NewImdbws(b))
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	for i := 0; i < b.N; i++ {
		if r, f := a.Parse(ctxTimeout, "tt0000014\tshort\tThe Waterer Watered\tL'arroseur arrosé\t0\t1895\t\\N\t1\tComedy,Short"); r != true || f != nil {
			b.Fatal("Unexpected response: false")
		}
	}
}

package analyzer

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"imdb.tsv.analyzer/pkg/imdbws"
)

func TestFilms_StoreAndGet(t *testing.T) {
	films := NewFilms()
	filmCh := make(chan imdbws.Film, 2)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		films.Store(ctx, filmCh)
	}()

	testFilms := []imdbws.Film{
		{Title: "Film 1"},
		{Title: "Film 2"},
	}

	filmCh <- testFilms[0]
	filmCh <- testFilms[1]
	close(filmCh)

	wg.Wait()

	storedFilms := films.GetFilms()
	if len(storedFilms) != len(testFilms) {
		t.Errorf("expected %d films, got %d", len(testFilms), len(storedFilms))
	}

	if !reflect.DeepEqual(testFilms, storedFilms) {
		t.Errorf("expected films %v, got %v", testFilms, storedFilms)
	}
}

func TestFilms_Store_ContextDone(t *testing.T) {
	films := NewFilms()
	filmCh := make(chan imdbws.Film, 2)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		films.Store(ctx, filmCh)
	}()

	film1 := imdbws.Film{Title: "Film 1"}
	filmCh <- film1

	// Give a moment for the first film to be processed
	time.Sleep(10 * time.Millisecond)

	cancel()

	// This second film should not be processed
	filmCh <- imdbws.Film{Title: "Film 2"}

	wg.Wait()

	storedFilms := films.GetFilms()
	if len(storedFilms) != 1 {
		t.Errorf("expected 1 film, got %d", len(storedFilms))
	}

	if storedFilms[0].Title != film1.Title {
		t.Errorf("expected film title %s, got %s", film1.Title, storedFilms[0].Title)
	}
}

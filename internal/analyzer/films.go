package analyzer

import (
	"context"
	"sync"

	"imdb.tsv.analyzer/internal/pkg/imdbws"
)

type Films struct {
	mx sync.Mutex

	films []imdbws.Film
}

func NewFilms() *Films {
	return &Films{}
}

func (f *Films) Store(ctx context.Context, film chan imdbws.Film) {
	for {
		select {
		case <-ctx.Done():
			return

		case l, ok := <-film:
			if !ok {
				return
			}

			f.mx.Lock()
			f.films = append(f.films, l)
			f.mx.Unlock()
		}
	}
}

func (f *Films) GetFilms() []imdbws.Film {
	f.mx.Lock()
	defer f.mx.Unlock()
	return f.films
}

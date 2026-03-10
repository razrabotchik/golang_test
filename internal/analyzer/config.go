package analyzer

import (
	env "github.com/caarlos0/env/v9"
	"imdb.tsv.analyzer/internal/pkg/imdbws"
)

type Config struct {
	MaxRunTime   uint `env:"MAX_RUN_TIME"`
	ThreadsCount uint `env:"THREADS_COUNT" envDefault:"10"`

	Imsdbws imdbws.Config `envPrefix:"IMDBWS_"`
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := env.ParseWithOptions(c, env.Options{}); err != nil {
		return nil, err
	}

	return c, nil
}

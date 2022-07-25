package infra

import (
	"flag"
	"os"
	"strconv"
)

type Options struct {
	ImdbApiKey     string
	FilePath       string
	Filters        Filters
	MaxApiRequests int // maximum number of requests to be made to [omdbapi](https://www.omdbapi.com/)
	MaxRunTime     int // maximum run time of the application. Format is a `time.Duration` string see [here](https://godoc.org/time#ParseDuration)
	MaxRequests    int // maximum number of requests to send to [omdbapi](https://www.omdbapi.com/)
	ThreadsCount   int
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	panic("the key " + key + " is not exists in the .env file")
}

func InitOptions() *Options {
	filters := Filters{}

	flag.StringVar(&filters.Id, "id", "", "filter id")
	flag.StringVar(&filters.TitleType, "titleType", "", "filter titleType")
	flag.StringVar(&filters.PrimaryTitle, "primaryTitle", "", "filter primaryTitle")
	flag.StringVar(&filters.OriginalTitle, "originalTitle", "", "filter originalTitle")
	flag.StringVar(&filters.Genre, "genre", "", "filter genre")
	flag.StringVar(&filters.StartYear, "startYear", "", "filter startYear")
	flag.StringVar(&filters.EndYear, "endYear", "", "filter endYear")
	flag.StringVar(&filters.RuntimeMinutes, "runtimeMinutes", "", "filter runtimeMinutes")
	flag.StringVar(&filters.Genres, "genres", "", "filter genres")
	flag.StringVar(&filters.PlotFilter, "plotFilter", "", "filter plotFilter")

	flag.Parse()

	var filePath string
	flag.StringVar(&filePath, "filePath", "", "absolute file path")
	if filePath == "" {
		filePath = getEnv("FILE_PATH")
	}

	threadsCount, _ := strconv.Atoi(getEnv("THREADS_COUNT"))
	if threadsCount <= 0 {
		threadsCount = 3
	}

	options := Options{
		FilePath:     filePath,
		Filters:      filters,
		ImdbApiKey:   getEnv("IMDB_API_KEY"),
		ThreadsCount: threadsCount,
	}

	flag.IntVar(&options.MaxApiRequests, "maxApiRequests", 0, "filter maxApiRequests")
	if options.MaxApiRequests == 0 {
		options.MaxApiRequests, _ = strconv.Atoi(getEnv("MAX_API_REQUESTS"))
	}

	flag.IntVar(&options.MaxRunTime, "maxRunTime", 0, "filter maxRunTime")
	if options.MaxRunTime == 0 {
		options.MaxRunTime, _ = strconv.Atoi(getEnv("MAX_RUN_TIME"))
	}
	flag.IntVar(&options.MaxRequests, "maxRequests", 0, "filter maxRequests")

	return &options
}

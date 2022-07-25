package infra

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type FilmRating struct {
	Source string
	Value  string
}

type Film struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Ratings    []FilmRating
	Metascore  string
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
}

type Imdbws struct {
	apiKey         string
	requestsCount  int
	maxApiRequests int
}

func NewImdbws(apiKey string, maxApiRequests int) Imdbws {
	return Imdbws{
		apiKey:         apiKey,
		requestsCount:  0,
		maxApiRequests: maxApiRequests,
	}
}

func (i *Imdbws) GetById(ctx context.Context, id string) *Film {
	i.requestsCount += 1

	if i.requestsCount > i.maxApiRequests {
		ShowErrorAndExit(errors.New("max API requests"))
		return nil
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://www.omdbapi.com/", nil)
	req.WithContext(ctx)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("apikey", i.apiKey)
	q.Add("i", id)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		ShowErrorAndExit(err)
		return nil
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var r *Film
	err = json.Unmarshal([]byte(respBody), &r)
	if err != nil {
		ShowErrorAndExit(err)
		return nil
	}

	return r
}

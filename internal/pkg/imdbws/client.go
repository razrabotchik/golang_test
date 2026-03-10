package imdbws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type response struct {
	Response ResponseResult `json:"Response"`
	Error    string         `json:"Error"`
}

func (r *response) IsError() bool {
	return r.Response == ResponseResultFalse
}

type ResponseResult string

const (
	ResponseResultTrue  ResponseResult = "True"
	ResponseResultFalse ResponseResult = "False"
)

type Client struct {
	cft           Config
	requestsCount atomic.Uint64

	http *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("apiKey is empty")
	}

	if cfg.ApiUrl == "" {
		return nil, errors.New("apiUrl is empty")
	}

	return &Client{
		cft:           cfg,
		requestsCount: atomic.Uint64{},
		http:          &http.Client{},
	}, nil
}

func (c *Client) GetFilmByID(ctx context.Context, ID string) (*Film, error) {
	requestsCount := c.requestsCount.Add(1)

	if c.cft.MaxApiRequests > 0 && requestsCount > uint64(c.cft.MaxApiRequests) {
		err := errors.New("max API requests")
		slog.ErrorContext(ctx, "max API requests", "err", err)
		return nil, err
	}

	req, _ := http.NewRequest("GET", c.cft.ApiUrl, nil)
	req = req.WithContext(ctx)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("apikey", c.cft.ApiKey)
	q.Add("i", ID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)

	if err != nil {
		slog.ErrorContext(ctx, "error imdbws get by id", "err", err)
		return nil, fmt.Errorf("error imdbws get by id: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, "error closing response body", "err", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "error imdbws read body", "err", err)
		return nil, fmt.Errorf("error imdbws read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.WarnContext(ctx, "error imdbws get by id: status code", "status", resp.StatusCode)

		apiResp := &response{}
		err = json.Unmarshal(respBody, apiResp)
		if err != nil {
			slog.ErrorContext(ctx, "error imdbws unmarshal", "err", err)
			return nil, fmt.Errorf("error imdbws unmarshal: %w", err)
		}

		if apiResp.IsError() {
			slog.ErrorContext(ctx, "error imdbws get by id", "err", apiResp.Error)
			return nil, fmt.Errorf("error imdbws get by id: %s", apiResp.Error)
		}

		return nil, fmt.Errorf("error imdbws get by id: %s", resp.Status)
	}

	var r *Film
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		slog.ErrorContext(ctx, "error imdbws unmarshal", "err", err)
		return nil, fmt.Errorf("error imdbws unmarshal: %w", err)
	}

	return r, nil
}

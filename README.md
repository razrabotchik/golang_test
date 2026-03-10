# IMDb TSV Analyzer

This project provides an example of processing TSV files and verifying data against the IMDb API.

The TSV file source is [https://datasets.imdbws.com/](https://datasets.imdbws.com/).

An API key can be generated from [https://www.omdbapi.com/apikey.aspx](https://www.omdbapi.com/apikey.aspx).

## Environment Variables

-   `IMDBWS_API_KEY`: Your IMDb API key.
-   `IMDBWS_API_URL`: The URL for the IMDb API.
-   `IMDBWS_MAX_API_REQUESTS`: The maximum number of requests to be made to [omdbapi](https://www.omdbapi.com/). Defaults to `1000` and can be overridden with the `--maxApiRequests` flag.
-   `MAX_RUN_TIME`: A timeout for the program's execution. No default.
-   `THREADS_COUNT`: The number of worker threads for file processing. Defaults to `10`.

## Usage

Place the TSV file on your local machine and use the `--file` flag to specify the path.

### Filters

The following filters are available as command-line flags:

-   `id`
-   `titleType`
-   `primaryTitle`
-   `originalTitle`
-   `genre`
-   `startYear`
-   `endYear`
-   `runtimeMinutes`
-   `genres`
-   `plotFilter`

### Example Command

Use the following command to run the program:

```shell
go run main.go --id=tt0000002 --titleType=short --file=data.tsv
```

## Testing
```shell
make test
```

## Benchmarking
```shell
make bench
```
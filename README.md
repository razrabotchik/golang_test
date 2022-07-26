Before runing:

<code>
cp .env.sample .env
</code>


Edit .env file and fill data for params:

- `FILE_PATH`
- `IMDB_API_KEY`

Put data.tsv. file to folder: ./data/data.tsv

There are filters

	- id
	- titleType
	- primaryTitle
	- originalTitle
	- genre
	- startYear
	- endYear
	- runtimeMinutes
	- genres
	- plotFilter 

To run use command:

<code>
go run main.go --id=tt0000002 --titleType=short
</code>


Tests:

<code>make tests</code>

Benchmark:

<code>make bench</code>

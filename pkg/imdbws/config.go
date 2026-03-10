package imdbws

type Config struct {
	ApiKey         string `env:"API_KEY,required"`
	ApiUrl         string `env:"API_URL" envDefault:"https://www.imdbws.com/"` // https://www.omdbapi.com/apikey.aspx
	MaxApiRequests uint   `env:"MAX_API_REQUESTS" default:"1000"`              // maximum number of requests to be made to [omdbapi](https://www.omdbapi.com/)
}

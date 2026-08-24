package serp

import (
	"github.com/serpapi/serpapi-golang"
	"os"
)

var client serpapi.SerpApiClient

func init() {
	// replace with your SerpApi key
	setting := serpapi.NewSerpApiClientSetting(os.Getenv("SERPAPI_KEY"))
	setting.Engine = "google_flights"
	// initialize the client
	client = serpapi.NewClient(setting)
}

// GetClient returns the configured client
func GetClient() serpapi.SerpApiClient {
	return client
}
